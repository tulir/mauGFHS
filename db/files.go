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
	"path"
)

// File represents a file ID to name link.
type File struct {
	ID                 string
	Name               string
	Namespace          string
	MIME               string
	DefaultPermissions PermissionValue
	namespace          *Namespace
	permissions        []Permission
}

const filesSchema = `
	id                 CHAR(32)          PRIMARY KEY,
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
	row := db.QueryRow(`SELECT id,name,namespace,mime,defaultPermissions FROM files WHERE id=?`, id)
	if row != nil {
		return scanFile(row)
	}
	return nil
}

// GetFileByPath gets a file by its namespace and name.
func GetFileByPath(namespace, name string) *File {
	row := db.QueryRow(`SELECT id,name,namespace,mime,defaultPermissions FROM files WHERE namespace=? AND name=?`, namespace, name)
	if row != nil {
		return scanFile(row)
	}
	return nil
}

func scanFile(row *sql.Row) *File {
	var id, name, namespace, mime string
	var defaultPermissions uint8
	row.Scan(&id, &name, &namespace, &mime, &defaultPermissions)
	return &File{ID: id, Name: name, Namespace: namespace, MIME: mime, DefaultPermissions: PermissionValue(defaultPermissions)}
}

// GetPermissions returns the permissions to this file.
func (file *File) GetPermissions() []Permission {
	if file.permissions != nil {
		return file.permissions
	}
	results, err := db.Query(`SELECT * FROM permissions WHERE file=?`, file.ID)
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
	return ioutil.ReadFile(path.Join(dataPath, file.ID))
}

// Path gets the display path of the file.
func (file *File) Path() string {
	return file.Namespace + "/" + file.Name
}
