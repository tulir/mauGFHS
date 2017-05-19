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
	"maunium.net/go/mauGFHS/db"
)

var config = configpkg.MainConfig

// Open opens the HTTP server.
func Open() {
	mainRouter := mux.NewRouter()
	r := mainRouter.PathPrefix(config.Listen.PathPrefix).Subrouter()
	r.Methods(http.MethodGet).Path("/direct/{id:[a-zA-Z0-9]{32}}").HandlerFunc(GetFileByID)
	r.Methods(http.MethodGet).Path("/direct/{namespace:[a-zA-Z0-9\\/]+}/{name}").HandlerFunc(GetFileByPath)

	server := &http.Server{
		Handler:      mainRouter,
		Addr:         fmt.Sprintf("%s:%d", config.Listen.Address, config.Listen.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatalln(server.ListenAndServe())
}

// GetFileByID handles an ID-based GET request.
func GetFileByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := db.GetFileByID(vars["id"])
	getFile(w, r, file)
}

// GetFileByPath handles a path-based GET request.
func GetFileByPath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := db.GetFileByPath(vars["name"], vars["namespace"])
	getFile(w, r, file)
}

func getFile(w http.ResponseWriter, r *http.Request, file *db.File) {
	if file == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := file.Read()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("Failed to read file %s: %v\n", file.Path(), err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
