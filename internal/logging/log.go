package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func Init(logDir string) error {
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return err
	}

	logFileName := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")

	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	multiOutInfo := io.MultiWriter(os.Stdout, file)
	multiOutError := io.MultiWriter(os.Stderr, file)

	Info = log.New(multiOutInfo, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(multiOutError, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
