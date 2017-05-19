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

import "database/sql"

// File represents a file ID to name link.
type File struct {
	ID        string
	Name      string
	Namespace string
	MIME      string
}

const filesSchema = `
	id CHAR(32) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	namespace VARCHAR(255) NOT NULL,
	mime VARCHAR(255) NOT NULL,
	UNIQUE KEY (name, namespace),
	CONSTRAINT namespace_name
		FOREIGN KEY (namespace) REFERENCES namespaces (name)
		ON DELETE CASCADE
		ON UPDATE RESTRICT
`

// GetFileByID gets a file by its storage ID.
func GetFileByID(id string) *File {
	row := db.QueryRow(`SELECT * FROM files WHERE id=?`, id)
	if row != nil {
		return scanFile(row)
	}
	return nil
}

// GetFileByPath gets a file by its namespace and name.
func GetFileByPath(name, namespace string) *File {
	row := db.QueryRow(`SELECT * FROM files WHERE name=? AND namespace=?`, name, namespace)
	if row != nil {
		return scanFile(row)
	}
	return nil
}

func scanFile(row *sql.Row) *File {
	var id, name, namespace, mime string
	row.Scan(&id, &name, &namespace, &mime)
	return &File{ID: id, Name: name, Namespace: namespace, MIME: mime}
}

// GetPermissions returns the permissions to this file.
func (file *File) GetPermissions() []Permission {
	results, err := db.Query(`SELECT * FROM permissions WHERE file=?`, file.ID)
	if err != nil {
		return []Permission{}
	}
	return scanFilePermissions(results)
}
