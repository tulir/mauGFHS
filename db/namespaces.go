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
	"strings"
)

// Namespace contains the details of a namespace.
type Namespace struct {
	Name      string
	Parent    string
	MIMETypes []string
}

const namespacesSchema = `
	name VARCHAR(255) PRIMARY KEY,
	parent VARCHAR(255) NOT NULL,
	mimes TEXT NOT NULL
`

func scanNamespace(row *sql.Row) *Namespace {
	var name, mimes string
	row.Scan(&name, &mimes)
	return &Namespace{Name: name, MIMETypes: strings.Split(mimes, ",")}
}

func scanNamespaces(results *sql.Rows) []*Namespace {
	data := []*Namespace{}
	for results.Next() {
		var name, parent, mimes string
		results.Scan(&name, &parent, &mimes)
		data = append(data, &Namespace{Name: name, MIMETypes: strings.Split(mimes, ",")})
	}
	return data
}

// GetNamespace gets the namespace with the given name from the database.
func GetNamespace(name string) *Namespace {
	row := db.QueryRow(`SELECT name,parent,mimes FROM namespaces WHERE name=?`, name)
	if row != nil {
		return scanNamespace(row)
	}
	return nil
}

// GetParent gets the parent of this namespace, or nil if this namespace doesn't have parent.
func (ns *Namespace) GetParent() *Namespace {
	if len(ns.Parent) == 0 {
		return nil
	}
	return GetNamespace(ns.Parent)
}

// GetChildren gets the namespaces that are children of this namespace.
func (ns *Namespace) GetChildren() []*Namespace {
	results, err := db.Query(`SELECT name,parent,mimes FROM namespaces WHERE parent=?`, ns.Name)
	if err == nil {
		return scanNamespaces(results)
	}
	return nil
}

// MIMETypesString turns the allowed MIME types array into a string.
func (ns *Namespace) MIMETypesString() string {
	return strings.Join(ns.MIMETypes, ",")
}

// Delete deletes this namespace from the database.
func (ns *Namespace) Delete() {
	db.Exec("DELETE FROM namespaces WHERE name=?", ns.Name)
}

// Update updates the database row for this namespace.
func (ns *Namespace) Update() {
	db.Exec("UPDATE namespaces SET mimetypes=? WHERE name=?", ns.MIMETypesString())
}

// Insert inserts this namespace definition into the database.
func (ns *Namespace) Insert() {
	db.Exec("INSERT INTO namespaces (name, mimetypes) VALUES (?, ?)", ns.Name, ns.MIMETypesString())
}
