package db

import "fmt"

// PermissionTargetType is the type of a permission target object.
type PermissionTargetType int

// The possible values for PermissionTargetType.
const (
	TypeFilePermission PermissionTargetType = iota
	TypeNamespacePermission
)

// PermissionValue is a int to permission enum mapping
type PermissionValue uint8

// Possible PermissionValues
const (
	PermissionNothing   PermissionValue = 0
	PermissionRead      PermissionValue = 1
	PermissionWrite     PermissionValue = 2
	PermissionReadWrite PermissionValue = PermissionRead + PermissionWrite
)

// CanRead checks if this PermissionValue is sufficient for reading files.
func (pv PermissionValue) CanRead() bool {
	return pv&PermissionRead == 1
}

// CanWrite checks if this PermissionValue is sufficient for writing files.
func (pv PermissionValue) CanWrite() bool {
	return pv&PermissionWrite == 1
}

// Permission is an abstract permission.
type Permission interface {
	GetUser() string
	GetTarget() string
	GetTargetType() PermissionTargetType
	GetPermission() PermissionValue
	SetPermission(pv PermissionValue)
	Delete()
	Insert()
	Update()
}

// UserPermissionsToMap turns a Permission array into a target -> permission map. This function
// completely ignores the user, see FilePermissionsToMap() for user -> permission mapping.
func UserPermissionsToMap(permissions []Permission) (data map[string]PermissionValue) {
	for _, permission := range permissions {
		data[permission.GetTarget()] = permission.GetPermission()
	}
	return
}

// TargetPermissionsToMap turns a Permission array into a user -> permission map. This function
// completely ignores the file, see UserPermissionsToMap() for target -> permission mapping.
func TargetPermissionsToMap(permissions []Permission) (data map[string]PermissionValue) {
	for _, permission := range permissions {
		data[permission.GetUser()] = permission.GetPermission()
	}
	return
}

type basePermission struct {
	User       string
	Target     string
	Permission PermissionValue
}

// GetTarget gets the target object of this permission.
func (perm *basePermission) GetTarget() string {
	return perm.Target
}

// GetUser gets the user that has this permission.
func (perm *basePermission) GetUser() string {
	return perm.User
}

// GetPermission gets the permission value in this key.
func (perm *basePermission) GetPermission() PermissionValue {
	return perm.Permission
}

// SetPermission gets the permission value in this key.
func (perm *basePermission) SetPermission(pv PermissionValue) {
	perm.Permission = pv
}

// Delete deletes this permission entry from the database.
func (perm *basePermission) Delete(tableName, targetFieldName string) {
	db.Exec(fmt.Sprintf("DELETE FROM %s WHERE user=? AND %s=?", tableName, targetFieldName), perm.User, perm.Target)
}

// Insert inserts this permission entry into the database.
func (perm *basePermission) Insert(tableName, targetFieldName string) {
	db.Exec(fmt.Sprintf("INSERT INTO %s (user, %s, permissions) VALUES (?, ?, ?)", tableName, targetFieldName), perm.User, perm.Target, perm.Permission)
}

// Update updates the permission value of this entry in the database.
func (perm *basePermission) Update(tableName, targetFieldName string) {
	db.Exec(fmt.Sprintf("UPDATE %s SET permissions=? WHERE user=? AND %s=?", tableName, targetFieldName), perm.Permission, perm.User, perm.Target)
}
