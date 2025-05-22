package logging

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func Init(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	logFile, err := os.OpenFile(dir+"/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	Info = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}
