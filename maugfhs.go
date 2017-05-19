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

package main

import (
	"fmt"
	"os"

	configpkg "maunium.net/go/mauGFHS/config"
	"maunium.net/go/mauGFHS/db"
	flag "maunium.net/go/mauflag"
	log "maunium.net/go/maulogger"
)

var configPath = flag.MakeFull("c", "config", "The path to the config file", "config.yml").String()
var debug = flag.MakeFull("d", "debug", "Whether or not to enable debug mode", "false").Bool()
var wantHelp, _ = flag.MakeHelpFlag()
var config = configpkg.MainConfig

func main() {
	flag.SetHelpTitles("mauGFHS 0.1 - A server that can serve as a backend for many kinds of services that only require file hosting.", "mauGFHS [-c /path/to/config] [-d] [-h]")
	err := flag.Parse()
	if err != nil || *wantHelp {
		flag.PrintHelp()
		return
	}

	err = configpkg.Open(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open config: %v\n", err)
		if *debug {
			panic(err)
		}
		os.Exit(1)
	}

	config.Logging.Configure(log.DefaultLogger)
	if *debug {
		log.DefaultLogger.PrintLevel = log.LevelDebug.Severity
	}
	log.Debugln("Logging initialized.")

	err = db.Open(config.Database)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v\n", err)
		if *debug {
			panic(err)
		}
		os.Exit(1)
	}
	db.CreateTables()

}
