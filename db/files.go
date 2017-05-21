// mauGFHS - A server that can serve as a backend for many kinds of services that only require file hosting.
// Copyright (C) 2017 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package db

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"path"

	log "maunium.net/go/maulogger"
)

// File represents a file ID to name link.
type File struct {
	ID                 string
	Size               int
	Name               string
	Namespace          string
	MIME               string
	DefaultPermissions PermissionValue
	namespace          *Namespace
	permissions        []Permission
}

const filesSchema = `
	id                 CHAR(32)          PRIMARY KEY,
	size               INTEGER           NOT NULL,
	name               VARCHAR(255)      NOT NULL,
	namespace          VARCHAR(255)      NOT NULL,
	mime               VARCHAR(255)      NOT NULL,
	defaultPermissions SMALLINT UNSIGNED NOT NULL,
	UNIQUE KEY (name, namespace),
	CONSTRAINT namespace_name
		FOREIGN KEY (namespace) REFERENCES namespaces (name)
		ON DELETE CASCADE
		ON UPDATE RESTRICT
`

// GetFileByID gets a file by its storage ID.
func GetFileByID(id string) *File {
	row := db.QueryRow(`SELECT id,size,name,namespace,mime,defaultPermissions FROM files WHERE id=?`, id)
	if row != nil {
		return scanFile(row)
	}
	return nil
}

// GetFileByPath gets a file by its namespace and name.
func GetFileByPath(namespace, name string) *File {
	row := db.QueryRow(`SELECT id,size,name,namespace,mime,defaultPermissions FROM files WHERE namespace=? AND name=?`, namespace, name)
	if row != nil {
		return scanFile(row)
	}
	return nil
}

func scanFile(row *sql.Row) *File {
	var id, name, namespace, mime string
	var size int
	var defaultPermissions uint8
	row.Scan(&id, &size, &name, &namespace, &mime, &defaultPermissions)
	return &File{ID: id, Size: size, Name: name, Namespace: namespace, MIME: mime, DefaultPermissions: PermissionValue(defaultPermissions)}
}

// Insert inserts this File into the database.
func (file *File) Insert() error {
	_, err := db.Exec(
		"INSERT INTO files (id,size,name,namespace,mime,defaultPermissions) VALUES (?, ?, ?, ?, ?, ?)",
		file.ID, file.Size, file.Name, file.Namespace, file.MIME, uint8(file.DefaultPermissions))
	return err
}

// SetDefaultPermissions sets the default permissions to this file.
func (file *File) SetDefaultPermissions(defaultPermissions PermissionValue) {
	file.DefaultPermissions = defaultPermissions
	db.Exec("UPDATE files SET defaultPermissions=%s WHERE id=%s", uint8(file.DefaultPermissions), file.ID)
}

// Delete deletes this file in the database.
func (file *File) Delete() {
	db.Exec("DELETE FROM files WHERE id=%s", file.ID)
}

// Rename changes the name of this File.
func (file *File) Rename(name string) error {
	_, err := db.Exec("UPDATE files SET name=%s WHERE id=%s", name, file.ID)
	if err != nil {
		return err
	}
	file.Name = name
	return nil
}

// Move moves this file into another namespace.
func (file *File) Move(namespace string) error {
	_, err := db.Exec("UPDATE files SET namespace=%s WHERE id=%s", namespace, file.ID)
	if err != nil {
		return err
	}
	file.Namespace = namespace
	return nil
}

// GetPermissionsFor gets the permissions to this file for a certain user. If the user is nil, the
// default permissions to the file will be returned.
func (file *File) GetPermissionsFor(user *User) PermissionValue {
	if user != nil {
		return user.GetPermissionValueToFile(file)
	}
	return file.DefaultPermissions
}

// GetPermissions returns the permissions to this file.
func (file *File) GetPermissions() []Permission {
	if file.permissions != nil {
		return file.permissions
	}
	results, err := db.Query(`SELECT user,file,permission FROM filepermissions WHERE file=?`, file.ID)
	if err != nil {
		return []Permission{}
	}
	file.permissions = scanFilePermissions(results)
	return file.permissions
}

// GetNamespace returns the namespace this file is in.
func (file *File) GetNamespace() *Namespace {
	if file.namespace == nil || file.namespace.Name != file.Namespace {
		file.namespace = GetNamespace(file.Namespace)
	}
	return file.namespace
}

// Read reads the file from disk.
func (file *File) Read() ([]byte, error) {
	data, err := ioutil.ReadFile(path.Join(dataPath, file.ID))
	if len(data) != file.Size {
		log.Warnf("File %s/%s had an unexpected size on disk! Expected: %d, got: %d", file.Namespace, file.Name, file.Size, len(data))
		file.Size = len(data)
		db.Exec("UPDATE files SET size=%s WHERE id=%s", file.Size, file.ID)
	}
	return data, err
}

// Write writes data for this file to the disk.
func (file *File) Write(data []byte) error {
	file.Size = len(data)
	file.MIME = http.DetectContentType(data)
	db.Exec("UPDATE files SET size=%s,mime=%s WHERE id=%s", file.Size, file.MIME, file.ID)
	return ioutil.WriteFile(path.Join(dataPath, file.ID), data, 0644)
}

// Path gets the display path of the file.
func (file *File) Path() string {
	return file.Namespace + "/" + file.Name
}
