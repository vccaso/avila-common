package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vccaso/avila-common/model"
	"github.com/vccaso/avila-common/util"
)

type KeyError struct{}

// MiddlewareValidationUser validates the error in the request and calls next if ok
func MiddlewareValidatorError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		error := model.Error{}
		util.CheckError(error.FromJson(req.Body))
		err := error.Validate()
		if err != nil {
			http.Error(rw, fmt.Sprintf("Error validating error: %v", err), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(req.Context(), KeyError{}, error)
		req = req.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
