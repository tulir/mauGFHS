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
	"io/ioutil"
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
	r.Methods(http.MethodGet).Path("/file/direct/{id:[a-zA-Z0-9]{32}}").HandlerFunc(GetFileByID)
	r.Methods(http.MethodGet).Path("/file/{namespace:[a-zA-Z0-9\\/]+}/{name}").HandlerFunc(GetFileByPath)
	r.Methods(http.MethodPut).Path("/file/direct/{id:[a-zA-Z0-9]{32}}").HandlerFunc(UpdateFileByID)
	r.Methods(http.MethodPut).Path("/file/{namespace:[a-zA-Z0-9\\/]+}/{name}").HandlerFunc(UpdateFileByPath)

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
	getFile(w, r, db.GetFileByID(vars["id"]))
}

// GetFileByPath handles a path-based GET request.
func GetFileByPath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	getFile(w, r, db.GetFileByPath(vars["name"], vars["namespace"]))
}

// UpdateFileByID handles an ID-based PUT request.
func UpdateFileByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	updateFile(w, r, db.GetFileByID(vars["id"]))
}

// UpdateFileByPath handles a path-based PUT request.
func UpdateFileByPath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	updateFile(w, r, db.GetFileByPath(vars["name"], vars["namespace"]))
}

func updateFile(w http.ResponseWriter, r *http.Request, file *db.File) {
	if file == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user := CheckAuth(r)
	if !file.GetPermissionsFor(user).CanWrite() {
		w.WriteHeader(403)
		return
	}

	r.ParseMultipartForm(32 << 20)

	fileData, _, err := r.FormFile("upload")
	if err != nil {
		w.WriteHeader(400)
		return
	}

	data, err := ioutil.ReadAll(fileData)
	if err != nil {
		log.Errorln("Failed to read data in form file!")
		w.WriteHeader(500)
		return
	}

	mime := http.DetectContentType(data)
	if !file.GetNamespace().IsMIMEAllowed(mime) {
		w.WriteHeader(415)
		return
	}

	err = file.Write(data, mime)
	if err != nil {
		log.Errorln("Failed to write file!")
		w.WriteHeader(500)
		return
	}
}

func getFile(w http.ResponseWriter, r *http.Request, file *db.File) {
	if file == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user := CheckAuth(r)
	if !file.GetPermissionsFor(user).CanRead() {
		w.WriteHeader(403)
		return
	}
	data, err := file.Read()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("Failed to read file %s: %v\n", file.Path(), err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
