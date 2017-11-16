// gosworld project main.go
package main

import (
	"flag"
	"log"
	"os"

	"regexp"

	"fmt"

	"github.com/kpawlik/smallgopher/worker"
)

const (
	version = "0.9"
)

//
// processName - name used to establish connection with SW ACP
// serverType - decide which server will be started
// configFilePath - path to config file
var (
	processName string
	port        string
	logFile     string
	offline     bool
)

// Init and parse command line params
func init() {
	var (
		file *os.File
		err  error
	)
	// declare and parse program flags
	flag.StringVar(&processName, "n", "", "process name")
	flag.StringVar(&port, "p", "", "Port number for worker eg \":4001\" ")
	flag.StringVar(&logFile, "l", "", "logfile")
	flag.BoolVar(&offline, "o", false, "Work with no Smallworld connection, only GeoJSON files")
	flag.Parse()
	if logFile == "" {
		log.SetOutput(os.Stdout)
	} else {
		if file, err = os.Create(logFile); err != nil {
			panic(err)
		}
		log.SetOutput(file)
	}
	if port == "" {
		log.Println("Missing port number")
		os.Exit(1)
	}
	if ok, _ := regexp.Match("^\\d+$", []byte(port)); ok {
		port = fmt.Sprintf(":%s", port)
	} else {
		if ok, _ := regexp.Match("^:\\d+$", []byte(port)); !ok {
			log.Printf("Wrong port number format %s. Expected eg. :4000 or 4000", port)
			os.Exit(1)
		}
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Version:     %s\n", version)
	log.Printf("Port: %s\n", port)
}

func main() {
	worker.StartWorker(port, processName, offline)
}
