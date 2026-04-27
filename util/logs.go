package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var logs_host string = "http://" + os.Getenv("LOGS_HOST") + "/api"

type ApiError struct {
	StatusCode int

	Message string
}

func HandleApiError(e error, appName string, httpCode int, errorMessage string, response http.ResponseWriter) {
	if e != nil {
		SendToLog(appName, "Error", e.Error())
	}
	if response != nil {
		response.WriteHeader(httpCode)
		jsonError, _ := json.Marshal(&ApiError{httpCode, errorMessage})
		fmt.Fprintf(response, "%v", string(jsonError))
	}
}
