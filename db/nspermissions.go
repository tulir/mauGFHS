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

// NamespacePermission contains the permissions that an user has to a namespace.
type NamespacePermission struct {
	basePermission
}

const nsPermissionsSchema = `
	user VARCHAR(255) NOT NULL,
	namespace VARCHAR(255) NOT NULL,
	permission SMALLINT UNSIGNED NOT NULL,
	PRIMARY KEY(user, file),
	CONSTRAINT user_email
		FOREIGN KEY (user) REFERENCES users (email)
		ON DELETE CASCADE
		ON UPDATE RESTRICT,
	CONSTRAINT namespace_name
		FOREIGN KEY (namespace) REFERENCES namespaces (name)
		ON DELETE CASCADE
		ON UPDATE RESTRICT
`

// GetTargetType gets the type of this permissions target object.
func (perm *NamespacePermission) GetTargetType() PermissionTargetType {
	return TypeNamespacePermission
}

// Delete deletes this permission entry from the database.
func (perm *NamespacePermission) Delete() {
	perm.basePermission.Delete("nspermissions", "namespace")
}

// Insert inserts this permission entry into the database.
func (perm *NamespacePermission) Insert() {
	perm.basePermission.Insert("nspermissions", "namespace")
}

// Update updates the permission value of this entry in the database.
func (perm *NamespacePermission) Update() {
	perm.basePermission.Update("nspermissions", "namespace")
}

func scanNamespacePermission(row *sql.Row) Permission {
	var user, namespace string
	var permission uint8
	row.Scan(&user, &namespace, &permission)
	return &NamespacePermission{basePermission{User: user, Target: namespace, Permission: PermissionValue(permission)}}
}

func scanNamespacePermissions(results *sql.Rows) []Permission {
	data := []Permission{}
	for results.Next() {
		var user, namespace string
		var permission uint8
		results.Scan(&user, &namespace, &permission)
		data = append(data, &NamespacePermission{basePermission{User: user, Target: namespace, Permission: PermissionValue(permission)}})
	}
	return data
}
