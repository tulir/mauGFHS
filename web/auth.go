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

package web

import (
	"net/http"

	"github.com/gorilla/sessions"

	"maunium.net/go/mauGFHS/db"
	log "maunium.net/go/maulogger"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

// CheckAuth checks if the given request is authenticated.
func CheckAuth(r *http.Request) *db.User {
	tokenStr := r.Header.Get("AuthToken")
	userStr := r.Header.Get("AuthUser")
	if len(tokenStr) > 0 && len(userStr) > 0 {
		user := db.GetUser(userStr)
		if user.CheckAuthToken(tokenStr) {
			return user
		}
		return nil
	}

	session, err := store.Get(r, "maugfhs")
	if err != nil {
		log.Errorln("Failed to check auth token cookie:", err)
		return nil
	}

	tokenStr = session.Values["authToken"].(string)
	userStr = session.Values["authUser"].(string)
	if len(tokenStr) > 0 && len(userStr) > 0 {
		user := db.GetUser(userStr)
		if user.CheckAuthToken(tokenStr) {
			return user
		}
		return nil
	}
	return nil
}
