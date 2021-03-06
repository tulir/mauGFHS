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
	"golang.org/x/crypto/bcrypt"
)

// User represents a row in the users table.
type User struct {
	Email    string
	Password []byte
	Admin    bool
}

const usersSchema = `
	email VARCHAR(255) PRIMARY KEY,
	password BINARY(60) NOT NULL,
	admin BOOLEAN NOT NULL
`

// GetUser gets the user with the given email.
func GetUser(email string) *User {
	row := db.QueryRow(`SELECT email,password,admin FROM users WHERE email=?`, email)
	if row != nil {
		var email string
		var password []byte
		var admin bool
		row.Scan(&email, &password, &admin)
		return &User{Email: email, Password: password, Admin: admin}
	}
	return nil
}

// CheckPassword checks if the given password is correct.
func (user *User) CheckPassword(password []byte) bool {
	return bcrypt.CompareHashAndPassword(user.Password, password) != nil
}

// ResetPassword resets the password of this user.
func (user *User) ResetPassword(newPassword []byte) bool {
	var err error
	user.Password, err = bcrypt.GenerateFromPassword(newPassword, bcrypt.DefaultCost)
	return err == nil
}

// GetAuthTokens gets the auth tokens of the user.
func (user *User) GetAuthTokens() []AuthToken {
	results, err := db.Query(`SELECT user,token,createdBy,expiry,isRecovery FROM authtokens WHERE user=? AND is_recovery=0`, user.Email)
	if err != nil {
		return []AuthToken{}
	}
	return scanAuthTokens(results)
}

// GetRecoveryTokens gets the password recovery tokens of the user.
func (user *User) GetRecoveryTokens() []AuthToken {
	results, err := db.Query(`SELECT user,token,createdBy,expiry,isRecovery FROM authtokens WHERE user=? AND is_recovery=1`, user.Email)
	if err != nil {
		return []AuthToken{}
	}
	return scanAuthTokens(results)
}

// CheckAuthToken checks if the given authentication token is valid for this user.
func (user *User) CheckAuthToken(token string) bool {
	for _, at := range user.GetAuthTokens() {
		if at.Token == token {
			return true
		}
	}
	return false
}

// CheckRecoveryToken checks if the given recovery token is valid for this user.
func (user *User) CheckRecoveryToken(token string) bool {
	for _, at := range user.GetRecoveryTokens() {
		if at.Token == token {
			return true
		}
	}
	return false
}

// GetPermissionsToFiles returns the file permissions this user has.
func (user *User) GetPermissionsToFiles() []Permission {
	results, err := db.Query(`SELECT user,namespace,permission FROM filepermissions WHERE user=?`, user.Email)
	if err != nil {
		return []Permission{}
	}
	return scanFilePermissions(results)
}

// GetPermissionToFile gets the permission this user has to the given file.
func (user *User) GetPermissionToFile(file *File) Permission {
	row := db.QueryRow(`SELECT user,file,permission FROM filepermissions WHERE user=? AND file=?`, user.Email, file.ID)
	if row == nil {
		return &FilePermission{basePermission{User: user.Email, Target: file.ID, Permission: PermissionNothing}}
	}
	return scanFilePermission(row)
}

func (user *User) GetPermissionValueToFile(file *File) PermissionValue {
	return user.GetPermissionToFile(file).GetPermission()
}

// GetPermissionsToNamespaces returns the namespace permissions this user has.
func (user *User) GetPermissionsToNamespaces() []Permission {
	results, err := db.Query(`SELECT user,namespace,permission FROM nspermissions WHERE user=?`, user.Email)
	if err != nil {
		return []Permission{}
	}
	return scanNamespacePermissions(results)
}

// GetPermissionToNamespace gets the permission this user has to the given namespace.
func (user *User) GetPermissionToNamespace(ns *Namespace) Permission {
	row := db.QueryRow(`SELECT user,file,permission FROM filepermissions WHERE user=? AND file=?`, user.Email, ns.Name)
	if row == nil {
		return &NamespacePermission{basePermission{User: user.Email, Target: ns.Name, Permission: PermissionNothing}}
	}
	return scanNamespacePermission(row)
}

func (user *User) GetPermissionValueToNamespace(ns *Namespace) PermissionValue {
	return user.GetPermissionToNamespace(ns).GetPermission()
}
