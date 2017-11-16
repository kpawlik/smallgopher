package server

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"github.com/gorilla/mux"
	"github.com/kpawlik/smallgopher/config"
)

var (
	reqFuncNameMap = map[string]string{
		"feature":         "Worker.GetFeature",
		"features":        "Worker.GetFeatures",
		"dump_features":   "Worker.DumpFeatures",
		"search_features": "Worker.SearchFeatures",
	}
	testReqFuncNameMap = map[string]string{
		"feature":         "Worker.GetTestFeature",
		"features":        "Worker.GetTestFeatures",
		"search_features": "Worker.TestSearchFeatures",
	}
)

// Response interface for responses objects
type Response interface {
	GetError() error
	GetBody() interface{}
	SetError(error)
	SetBody(interface{})
}

// WorkerConnection type to store worker connection and data
type WorkerConnection struct {
	Name string
	Host string
	Port int
	Conn *rpc.Client
}

// channels with workers connections
type workerChan chan *WorkerConnection

// StartServer initialize workers and starts HTTP server
func StartServer(config *config.Config, offline bool) {
	onlineWorkers, offlineWorkers := initWorkersConnections(config.Workers)
	// try reconnect workers in goroutine
	go handleOfflineWorkers(onlineWorkers, offlineWorkers)
	startHTTPServer(config, onlineWorkers, offlineWorkers, offline)
	defer func() {
		if err := recover(); err != nil {
			log.Panic("PANIC ", err)
		}
	}()
}

// startHTTPServer starts main http server
func startHTTPServer(config *config.Config, onlineWorkers workerChan, offlineWorkers workerChan, offline bool) {
	port := config.Server.Port
	log.Printf("HTTP server started on port: %v\n", port)

	handler := &ReqHandler{OnlineWorkers: onlineWorkers,
		OfflineWorkers: offlineWorkers,
		Config:         config,
		Offline:        offline}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.ServeMain)
	r.HandleFunc("/config", handler.ServeConfig)
	r.HandleFunc("/config/{feature}", handler.ServeConfig)
	r.HandleFunc("/features/{feature}/{bbox}/{limit}", handler.ServeFeatures)
	r.HandleFunc("/dump_features/{feature}/{bbox}/{limit}", handler.DumpFeatures)
	r.HandleFunc("/feature/{feature}/{id}", handler.ServeFeature)
	r.HandleFunc("/search", handler.SearchFeature).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// initWorkersConnections initialize RPC clients connections. Cache them to online channel.
// Workers which are not connected are send to offline channel
func initWorkersConnections(workersDef []*config.WorkerConf) (onlineWorkers workerChan, offlineWorkers workerChan) {
	workersNo := len(workersDef)
	onlineWorkers = make(workerChan, workersNo)
	offlineWorkers = make(workerChan, workersNo)

	for _, workerDef := range workersDef {
		conn, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", workerDef.Host, workerDef.Port))
		if err != nil {
			// add to online pool
			offlineWorkers <- &WorkerConnection{Host: workerDef.Host,
				Port: workerDef.Port,
				Name: workerDef.Name}
			continue
		}
		// add to offline pool
		onlineWorkers <- &WorkerConnection{Host: workerDef.Host,
			Port: workerDef.Port,
			Name: workerDef.Name,
			Conn: conn}
		log.Printf("Worker %s (%s:%d), CONNECTED \n", workerDef.Name, workerDef.Host, workerDef.Port)
	}
	return
}

//handleOfflineWorkers trying to reconnect offline workers every one second. When worker will reconnect
// send him to online chanel
func handleOfflineWorkers(onlineWorkers workerChan, offlineWorkers workerChan) {
	for {
		<-time.After(1 * time.Second)
		for worker := range offlineWorkers {
			conn, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", worker.Host, worker.Port))
			if err != nil {
				offlineWorkers <- worker
			} else {
				log.Printf("Worker %s (%s:%d), CONNECTED \n", worker.Name, worker.Host, worker.Port)
				worker.Conn = conn
				onlineWorkers <- worker
			}
		}
	}
}
