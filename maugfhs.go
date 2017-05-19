package main

import (
	"fmt"
	"os"

	"maunium.net/go/mauGFHS/config"
	"maunium.net/go/mauGFHS/db"
	flag "maunium.net/go/mauflag"
	log "maunium.net/go/maulogger"
)

var configPath = flag.MakeFull("c", "config", "The path to the config file", "config.yml").String()
var debug = flag.MakeFull("d", "debug", "Whether or not to enable debug mode", "false").Bool()
var wantHelp, _ = flag.MakeHelpFlag()

func main() {
	flag.SetHelpTitles("mauGFHS 0.1 - A server that can serve as a backend for many kinds of services that only require file hosting.", "mauGFHS [-c /path/to/config] [-d] [-h]")
	err := flag.Parse()
	if err != nil || *wantHelp {
		flag.PrintHelp()
		return
	}

	err = config.Open(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open config: %v\n", err)
		if *debug {
			panic(err)
		}
		os.Exit(1)
	}

	config.GetConfig().Logging.Configure(log.DefaultLogger)
	if *debug {
		log.DefaultLogger.PrintLevel = log.LevelDebug.Severity
	}
	log.Debugln("Logging initialized.")

	err = db.Open(config.GetConfig().Database)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v\n", err)
		if *debug {
			panic(err)
		}
		os.Exit(1)
	}
	db.CreateTables()

}
