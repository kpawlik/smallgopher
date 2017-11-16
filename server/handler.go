package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kpawlik/smallgopher/config"
)

// ReqHandler struct to implement ServeHTTP method
type ReqHandler struct {
	OnlineWorkers  workerChan
	OfflineWorkers workerChan
	Config         *config.Config
	Offline        bool
}

// getRequestFunctionName return name of method which will handle RCP request.
// Second return value is bool, which indicates if function name was found or not.
func (r *ReqHandler) getRequestFunctionName(protocolName string) (string, bool) {
	var (
		funcMap map[string]string
	)
	if r.Offline {
		funcMap = testReqFuncNameMap
	} else {
		funcMap = reqFuncNameMap
	}
	funcName, ok := funcMap[protocolName]
	return funcName, ok
}

//writeErrorStatus writes error on writer, error message depends on status value
func (r *ReqHandler) writeErrorStatus(w http.ResponseWriter, status int) {
	var message string
	switch status {
	case http.StatusMethodNotAllowed:
		message = "Unsupported protocol"
	case http.StatusUnauthorized:
		message = "Unauthorized protocol"
	}
	w.WriteHeader(status)
	fmt.Fprintf(w, message)
}

// parseBBox convert string to bbox
func (r *ReqHandler) parseBBox(bboxString string) (bbox []float32, err error) {
	var (
		res []string
		fl  float64
	)
	if bboxString, err = url.QueryUnescape(bboxString); err != nil {
		log.Printf("Error unescape bbox: %v\n", err)
	}
	res = strings.Split(bboxString, ",")
	bbox = make([]float32, 4, 4)
	for i, s := range res {
		if fl, err = strconv.ParseFloat(s, 32); err != nil {
			return
		}
		bbox[i] = float32(fl)
	}
	return
}

func (r *ReqHandler) addJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

// ServeMain is http request handler.
func (r *ReqHandler) ServeMain(w http.ResponseWriter, req *http.Request) {
	data := struct {
		Center []float32
	}{
		Center: r.Config.Server.Center,
	}
	template.Must(template.ParseFiles("static/html/app.html")).Execute(w, data)
}

// ServeConfig is http request handler.
func (r *ReqHandler) ServeConfig(w http.ResponseWriter, req *http.Request) {
	var (
		featureConf interface{}
	)
	vars := mux.Vars(req)
	featureName := vars["feature"]
	// serve for feature with name or for all features
	if featureName != "" {
		featureConf = r.Config.GetFeaturesDef(featureName)
	} else {
		featureConf = r.Config.Server.Features
	}
	r.addJSONHeader(w)
	enc := json.NewEncoder(w)
	if err := enc.Encode(featureConf); err != nil {
		r.logReturnError(w, fmt.Sprintf("JSON encode config error: %v", err))
		return
	}
}

func (r *ReqHandler) sendRequest(requestFuncName string, request interface{}, w http.ResponseWriter) (response *config.FeaturesResponse, err error) {
	response = &config.FeaturesResponse{}
	for {
		// get free worker from online pool
		worker := <-r.OnlineWorkers
		conn := worker.Conn
		if err = conn.Call(requestFuncName, request, response); err == nil {
			// return worker to the online pool
			r.OnlineWorkers <- worker
			break
		} else {
			err = fmt.Errorf("Error: %s (Worker %s). ", err, worker.Name)
			r.logReturnError(w, err.Error())
			// add worker to the offline pool and get next worker from online pool. dont break
			hj, ok := w.(http.Hijacker)
			if !ok {
				return
			}
			httpConn, _, httpErr := hj.Hijack()
			if httpErr != nil {
				return
			}
			httpConn.Close()
			r.OfflineWorkers <- worker
		}
	}
	return
}

func (r *ReqHandler) encodeResponse(w http.ResponseWriter, response *config.FeaturesResponse) {
	var (
		err error
	)
	if err = response.GetError(); err != nil {
		r.logReturnError(w, fmt.Sprint(err))
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(response.GetBody()); err != nil {
		r.logReturnError(w, fmt.Sprintf("Error Encode response body: %s", err))
		return
	}
}

// ServeFeatures is http request handler.
func (r *ReqHandler) ServeFeatures(w http.ResponseWriter, req *http.Request) {
	var (
		requestFuncName string
		ok              bool
		err             error
	)
	defer func() {
		if err := recover(); err != nil {
			log.Panicf("PANIC: %v\n", err)
		}
	}()
	r.addJSONHeader(w)
	vars := mux.Vars(req)
	bbox, limit, featureConf, err := r.validateFeaturesParams(vars)
	if err != nil {
		r.logReturnError(w, fmt.Sprint(err))
		return
	}
	// check request function name
	if requestFuncName, ok = r.getRequestFunctionName("features"); !ok {
		r.writeErrorStatus(w, http.StatusMethodNotAllowed)
		return
	}
	request := config.NewFeaturesRequest(featureConf, bbox, limit)
	response, err := r.sendRequest(requestFuncName, request, w)
	if err != nil {
		r.logReturnError(w, fmt.Sprint(err))
		return
	}
	r.encodeResponse(w, response)
}

// ServeFeature is http request handler.
func (r *ReqHandler) ServeFeature(w http.ResponseWriter, req *http.Request) {
	var (
		requestFuncName string
		ok              bool
		err             error
	)
	defer func() {
		if err := recover(); err != nil {
			log.Panicf("PANIC: %v\n", err)
		}
	}()
	r.addJSONHeader(w)
	vars := mux.Vars(req)
	featureName := vars["feature"]
	id := vars["id"]
	featureConf := r.Config.GetFeaturesDef(featureName)
	// check request function name
	if requestFuncName, ok = r.getRequestFunctionName("feature"); !ok {
		r.writeErrorStatus(w, http.StatusMethodNotAllowed)
		return
	}
	request, err := config.NewFeatureRequest(featureConf, id)
	if err != nil {
		r.logReturnError(w, err.Error())
		return
	}
	response, err := r.sendRequest(requestFuncName, request, w)
	if err != nil {
		r.logReturnError(w, fmt.Sprint(err))
		return
	}
	r.encodeResponse(w, response)
}

// DumpFeatures is http request handler.
func (r *ReqHandler) DumpFeatures(w http.ResponseWriter, req *http.Request) {
	var (
		requestFuncName string
		ok              bool
		err             error
	)
	defer func() {
		if err := recover(); err != nil {
			log.Panicf("PANIC: %v\n", err)
		}
	}()
	r.addJSONHeader(w)
	vars := mux.Vars(req)
	bbox, limit, featureConf, err := r.validateFeaturesParams(vars)
	if err != nil {
		r.logReturnError(w, fmt.Sprint(err))
		return
	}
	// check request function name
	if requestFuncName, ok = r.getRequestFunctionName("dump_features"); !ok {
		r.writeErrorStatus(w, http.StatusMethodNotAllowed)
		return
	}
	request := config.NewDumpRequest(featureConf, bbox, limit)
	response, err := r.sendRequest(requestFuncName, request, w)
	if err != nil {
		r.logReturnError(w, fmt.Sprint(err))
		return
	}
	r.encodeResponse(w, response)

}

// SearchFeature is http request handler.
func (r *ReqHandler) SearchFeature(w http.ResponseWriter, req *http.Request) {
	var (
		requestFuncName string
		ok              bool
		err             error
		body            []byte
		response        *config.FeaturesResponse
	)
	defer func() {
		if err := recover(); err != nil {
			log.Panicf("PANIC: %v\n", err)
		}
	}()
	r.addJSONHeader(w)
	if body, err = ioutil.ReadAll(req.Body); err != nil {
		r.logReturnError(w, err)
		return
	}
	searchParams := struct {
		Name      string   `json:"name"`
		FieldName string   `json:"fieldName"`
		Values    []string `json:"values"`
	}{}
	if err = json.Unmarshal(body, &searchParams); err != nil {
		r.logReturnError(w, err)
		return
	}
	featureConf := r.Config.GetFeaturesDef(searchParams.Name)
	fieldConf := featureConf.GetField(searchParams.FieldName)
	if requestFuncName, ok = r.getRequestFunctionName("search_features"); !ok {
		r.writeErrorStatus(w, http.StatusMethodNotAllowed)
		return
	}
	request := config.NewSearchRequest(featureConf, fieldConf.Name, fieldConf.Type, searchParams.Values)
	if response, err = r.sendRequest(requestFuncName, request, w); err != nil {
		r.logReturnError(w, err)
		return
	}
	r.encodeResponse(w, response)
}

func (r *ReqHandler) logReturnError(w http.ResponseWriter, err interface{}) {
	msg := ""
	switch err.(type) {
	case error:
		msg = fmt.Sprintf(`{"error": "%s"}`, err.(error))
	case string:
		msg = fmt.Sprintf(`{"error": "%s"}`, err.(string))
	default:
		msg = fmt.Sprintf(`{"error": "%v"}`, err)
	}
	fmt.Fprintf(w, msg)
	log.Printf(msg)
}

func (r *ReqHandler) validateFeaturesParams(vars map[string]string) (bbox []float32, limit int, conf *config.FeatureConf, err error) {
	bbox, err = r.parseBBox(vars["bbox"])
	if err != nil {
		err = fmt.Errorf("Error parse BBOX: %s", vars["bbox"])
		return
	}
	featureName := vars["feature"]
	limit, err = strconv.Atoi(vars["limit"])
	if err != nil {
		err = fmt.Errorf("Error parse limit: %s", vars["limit"])
		return
	}
	conf = r.Config.GetFeaturesDef(featureName)
	if conf == nil {
		err = fmt.Errorf("No config found for '%s'", vars["feature"])
		return
	}
	return
}
