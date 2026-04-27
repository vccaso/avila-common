package http

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/vccaso/avila-common/util"
)

// RecoveryMiddleware catches panics and logs them to avila-logs
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				appName := os.Getenv("APPLICATION_NAME")
				if appName == "" {
					appName = "unknown-app"
				}

				stackTrace := debug.Stack()
				errorMessage := fmt.Sprintf("PANIC: %v\nStack Trace:\n%s", err, stackTrace)
				
				// Log to console
				util.Error.Println(errorMessage)
				
				// Send to avila-logs
				util.SendToLog(appName, "PANIC", errorMessage)

				// Respond with 500
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"status": 500, "message": "Internal Server Error (Panic caught)"}`)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
