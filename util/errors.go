package util

import (
	"log"
	"os"
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

func CheckError(e error) {
	if e != nil {
		Error.Println(e.Error())
	}
}

func CheckWarning(e error) {
	if e != nil {
		Warning.Println(e.Error())
	}
}

func CheckInfo(e error) {
	if e != nil {
		Info.Println(e.Error())
	}
}

func CheckDebug(e error) {
	if e != nil {
		Debug.Println(e.Error())
	}
}
