package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vccaso/avila-common/model"
)

var logs_host string = "http://" + os.Getenv("LOGS_HOST") + "/api"

type ApiError struct {
	StatusCode int

	Message string
}

func HandleApiError(e error, appName string, httpCode int, errorMessage string, response http.ResponseWriter) {

	if e != nil {
		error := model.Error{}
		error.App = appName
		error.Error_time = time.Now()
		error.Level = "Error"
		error.Message = e.Error()
		error.Gateway_session = "TODO"

		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(error)
		resp, err := http.Post(logs_host+"/error", "application/json", bytes.NewBuffer(reqBodyBytes.Bytes()))
		CheckError(err)
		if resp != nil {
			defer resp.Body.Close()
		}
	}
	if response != nil {

		response.WriteHeader(httpCode)
		jsonError, _ := json.Marshal(&ApiError{httpCode, errorMessage})
		fmt.Fprintf(response, "%v", string(jsonError))
	}
}
