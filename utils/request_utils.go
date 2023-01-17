package utils

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

// Gets the valid path parameter for a request from the given http request
func GetPathParam(r *http.Request, paramName string) (paramValue string, err error) {
	paramVals := []string{paramName}
	params := GetPathParams(r, paramVals)

	pathParam, ok := params[paramName]
	paramValue, strOk := pathParam.(string)

	if !ok || !strOk {
		err = errors.New("error retrieving " + paramName + " from request path")
	}
	return
}

func GetPathParams(r *http.Request, paramVals []string) (params map[string]interface{}) {
	params = make(map[string]interface{})
	vars := mux.Vars(r)
	for _, param := range paramVals {
		value, ok := vars[param]
		if ok {
			params[param] = value
		}
	}
	return
}

func GetQueryParams(r *http.Request, paramVals []string) (params map[string]interface{}) {
	params = make(map[string]interface{})
	values := r.URL.Query()

	for _, param := range paramVals {
		value := values.Get(param)
		if len(value) > 0 {
			params[param] = value
		}
	}
	return
}

// RetrievePayload is an utility to convert payload based on the Content-Type header information set in the request
func RetrievePayload(r *http.Request) (interface{}, error) {

	// retrieve JSON/ XML payload and use appropriate utils to transform the information based on Content-Type header information
	switch r.Header.Get(ContentType) {

	case JSONContentType:
		return retrieveJSONPayload(r)
	case XMLContentType:
		return "", errors.New("XML is not a supported format at the moment, please change the header value to application/json")
	default:
		return retrieveJSONPayload(r)
	}
}

func retrieveJSONPayload(r *http.Request) (interface{}, error) {
	var payload interface{}

	if r.ContentLength > 0 {
		body, readError := ioutil.ReadAll(r.Body)
		if readError != nil {
			return nil, readError
		}

		jsonError := json.Unmarshal(body, &payload)
		if jsonError != nil {
			return nil, jsonError
		}
	}

	return payload, nil
}
