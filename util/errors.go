package util

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vccaso/avila-common/model"
)

// Info writes logs in the color blue with "INFO: " as prefix
var Info = log.New(os.Stdout, "\u001b[34mINFO: \u001B[0m", log.LstdFlags|log.Lshortfile)

// Warning writes logs in the color yellow with "WARNING: " as prefix
var Warning = log.New(os.Stdout, "\u001b[33mWARNING: \u001B[0m", log.LstdFlags|log.Lshortfile)

// Error writes logs in the color red with "ERROR: " as prefix
var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)

// Debug writes logs in the color cyan with "DEBUG: " as prefix
var Debug = log.New(os.Stdout, "\u001b[36mDEBUG: \u001B[0m", log.LstdFlags|log.Lshortfile)

func CheckPanic(e error) {
	if e != nil {
		panic(e.Error())
	}
}

var app string = os.Getenv("APPLICATION_NAME")

func CheckError(e error) {
	if e != nil {
		Error.Println(e.Error())
		SendToLog(app, "ERROR", e.Error())
	}
}

func LogError(e error) {
	if e != nil {
		Error.Println(e.Error())
	}
}

func CheckWarning(m string) {
	Warning.Println(m)
	SendToLog(app, "WARNING", m)
}

func CheckInfo(m string) {
	Info.Println(m)
	SendToLog(app, "INFO", m)
}

func CheckDebug(m string) {
	Debug.Println(m)
	SendToLog(app, "DEBUG", m)
}

func SendToLog(app string, level string, m string) {
	go func() {
		errorModel := model.Error{}
		errorModel.App = app
		errorModel.Error_time = time.Now()
		errorModel.Message = m
		errorModel.Gateway_session = "TODO"
		errorModel.Level = level

		reqBodyBytes := new(bytes.Buffer)
		if err := json.NewEncoder(reqBodyBytes).Encode(errorModel); err != nil {
			Error.Println("Error encoding log:", err)
			return
		}

		client := &http.Client{
			Timeout: 5 * time.Second,
		}
		req, err := http.NewRequest(http.MethodPost, logs_host+"/error", reqBodyBytes)
		if err != nil {
			Error.Println("Error creating log request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			Error.Println("Error sending log to avila-logs:", err)
			return
		}
		defer resp.Body.Close()
	}()
}
