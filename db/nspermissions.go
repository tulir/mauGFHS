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
