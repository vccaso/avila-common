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
	error := model.Error{}
	error.App = app
	error.Error_time = time.Now()
	error.Message = m
	error.Gateway_session = "TODO"
	error.Level = level

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(error)
	resp, err := http.Post(logs_host+"/error", "application/json", bytes.NewBuffer(reqBodyBytes.Bytes()))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		Error.Println(err.Error())
	}

}
