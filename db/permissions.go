package db

import "database/sql"

// Permission represents a row in the permissions table.
type Permission struct {
	Owner      string
	File       string
	Permission uint8
}

const permissionsSchema = `
	owner VARCHAR(255) NOT NULL,
	file CHAR(32) NOT NULL,
	permission SMALLINT UNSIGNED NOT NULL,
	PRIMARY KEY(owner, file),
	CONSTRAINT owner_email
		FOREIGN KEY (owner) REFERENCES users (email)
		ON DELETE CASCADE
		ON UPDATE RESTRICT,
	CONSTRAINT file_id
		FOREIGN KEY (file) REFERENCES files (id)
		ON DELETE CASCADE
		ON UPDATE RESTRICT
`

func scanPermission(row *sql.Row) *Permission {
	var owner, file string
	var permission uint8
	row.Scan(&owner, &file, &permission)
	return &Permission{Owner: owner, File: file, Permission: permission}
}

func scanPermissions(results *sql.Rows) []*Permission {
	data := []*Permission{}
	for results.Next() {
		var owner, file string
		var permission uint8
		results.Scan(&owner, &file, &permission)
		data = append(data, &Permission{Owner: owner, File: file, Permission: permission})
	}
	return data
}

// FilePermissionsToMap turns a Permission array into a user -> permission map. This function
// completely ignores the file, see UserPermissionsToMap() for file -> permission mapping.
func FilePermissionsToMap(permissions []*Permission) (data map[string]uint8) {
	for _, permission := range permissions {
		data[permission.Owner] = permission.Permission
	}
	return
}

// UserPermissionsToMap turns a Permission array into a file -> permission map. This function
// completely ignores the owner, see FilePermissionsToMap() for owner -> permission mapping.
func UserPermissionsToMap(permissions []*Permission) (data map[string]uint8) {
	for _, permission := range permissions {
		data[permission.File] = permission.Permission
	}
	return
}
