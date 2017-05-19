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
	"fmt"
	"net/http"
	"time"

	log "maunium.net/go/maulogger"

	"github.com/gorilla/mux"
	configpkg "maunium.net/go/mauGFHS/config"
)

var config = configpkg.MainConfig

// Open opens the HTTP server.
func Open() {
	mainRouter := mux.NewRouter()
	r := mainRouter.PathPrefix(config.Listen.PathPrefix).Subrouter()
	r.Methods("GET", "PUT", "DELETE").HandleFunc("/direct/{id:[a-zA-Z0-9]{32}}", nil)
	r.Methods("GET", "PUT", "POST", "DELETE").HandlerFunc("/", nil)

	server := &http.Server{
		Handler:      mainRouter,
		Addr:         fmt.Sprintf("%s:%d", config.Listen.Address, config.Listen.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatalln(server.ListenAndServe())
}
