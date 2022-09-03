package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Helper functions to write a JSON response to the client

func JsonResponseSuccess(response http.ResponseWriter, data interface{}) {

	response.WriteHeader(http.StatusOK)
	ljson, _ := json.Marshal(data)
	fmt.Fprintf(response, "%v", string(ljson))
}

func JsonResponseWithCode(response http.ResponseWriter, data interface{}, httpCode int) {

	response.WriteHeader(httpCode)
	ljson, _ := json.Marshal(data)
	fmt.Fprintf(response, "%v", string(ljson))
}

func ObjectFromBody(request *http.Request, data interface{}) error {

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&data)
	return err
}
