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
	"time"
)

// AuthToken represents either an authentication token or a password recovery token.
type AuthToken struct {
	User   string
	Token  string
	Expiry int64
}

const authTokensSchema = `
	user VARCHAR(255) NOT NULL,
	token VARCHAR(64) NOT NULL,
	createdBy VARCHAR(255) NOT NULL,
	expiry BIGINT NOT NULL,
	isRecovery BOOLEAN NOT NULL DEFAULT '0',
	PRIMARY KEY (user, token),
	CONSTRAINT user_email
		FOREIGN KEY (user) REFERENCES users (email)
		ON DELETE CASCADE
		ON UPDATE RESTRICT
`

// Delete this auth token from the database.
func (at AuthToken) Delete() {
	db.Exec("DELETE FROM authtokens WHERE user=? AND token=?", at.User, at.Token)
}

// HasExpired checks if the auth token has expired.
func (at AuthToken) HasExpired() bool {
	return at.Expiry < time.Now().Unix()
}

func scanAuthTokens(results *sql.Rows) []AuthToken {
	data := []AuthToken{}
	for results.Next() {
		var email, token, createdBy string
		var expiry int64
		var isRecovery bool
		results.Scan(&email, &token, &createdBy, &expiry, &isRecovery)
		at := AuthToken{User: email, Token: token, Expiry: expiry}
		if !at.HasExpired() {
			data = append(data, at)
		}
	}
	return data
}
