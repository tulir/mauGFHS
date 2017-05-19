package db

import "database/sql"

// FilePermission contains the permissions that an user has to a file.
type FilePermission struct {
	basePermission
}

const filePermissionsSchema = `
	user VARCHAR(255) NOT NULL,
	file CHAR(32) NOT NULL,
	permission SMALLINT UNSIGNED NOT NULL,
	PRIMARY KEY(user, file),
	CONSTRAINT user_email
		FOREIGN KEY (user) REFERENCES users (email)
		ON DELETE CASCADE
		ON UPDATE RESTRICT,
	CONSTRAINT file_id
		FOREIGN KEY (file) REFERENCES files (id)
		ON DELETE CASCADE
		ON UPDATE RESTRICT
`

// GetTargetType gets the type of this permissions target object.
func (perm *FilePermission) GetTargetType() PermissionTargetType {
	return TypeFilePermission
}

// Delete deletes this permission entry from the database.
func (perm *FilePermission) Delete() {
	perm.basePermission.Delete("filepermissions", "file")
}

// Insert inserts this permission entry into the database.
func (perm *FilePermission) Insert() {
	perm.basePermission.Insert("filepermissions", "file")
}

// Update updates the permission value of this entry in the database.
func (perm *FilePermission) Update() {
	perm.basePermission.Update("filepermissions", "file")
}

func scanFilePermission(row *sql.Row) Permission {
	var user, file string
	var permission uint8
	row.Scan(&user, &file, &permission)
	return &FilePermission{basePermission{User: user, Target: file, Permission: PermissionValue(permission)}}
}

func scanFilePermissions(results *sql.Rows) []Permission {
	data := []Permission{}
	for results.Next() {
		var user, file string
		var permission uint8
		results.Scan(&user, &file, &permission)
		data = append(data, &FilePermission{basePermission{User: user, Target: file, Permission: PermissionValue(permission)}})
	}
	return data
}
