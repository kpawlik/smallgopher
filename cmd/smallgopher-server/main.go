// gosworld project main.go
package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/kpawlik/smallgopher/config"
	"github.com/kpawlik/smallgopher/server"
)

const (
	version = "0.9"
)

//
// serverType - decide which server will be started
// configFilePath - path to config file
var (
	processName    string
	serverType     string
	configFilePath string
	logFile        string
	offline        bool
)

// Init and parse command line params
func init() {
	var (
		file *os.File
		err  error
	)

	runtime.GOMAXPROCS(runtime.NumCPU())
	// declare and parse program flags
	flag.StringVar(&configFilePath, "c", "", "path to config file")
	flag.StringVar(&logFile, "l", "", "logfile")
	flag.BoolVar(&offline, "o", false, "Work with no Smallworld connection, only GeoJSON files")
	flag.Parse()
	if configFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}
	if logFile == "" {
		log.SetOutput(os.Stdout)
	} else {
		if file, err = os.Create(logFile); err != nil {
			panic(err)
		}
		log.SetOutput(file)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Version:     %s\n", version)
	log.Printf("Config file: %s\n", configFilePath)
	if offline {
		log.Printf("Offline mode")
	}
}

func main() {
	var (
		cfg *config.Config
		err error
	)

	if cfg, err = config.ReadConf(configFilePath); err != nil {
		log.Panicf("Error reading config file: %v\n", err)
	}
	startHTTPServer(cfg, offline)
}

func startHTTPServer(config *config.Config, offline bool) {
	server.StartServer(config, offline)
}
