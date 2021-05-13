package logger

import (
	"io"
	"log"
	"os"
)

var (
	Info    *log.Logger // Important information
	Warning *log.Logger // Be concerned
	Error   *log.Logger // Critical problem
	Fatal   *log.Logger // Fatal error
)

//NewLogger sets logging parameter and the output file
func NewLogger(logFile string) (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return file, err
	}
	Info = log.New(file, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(file, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stderr), "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	Fatal = log.New(io.MultiWriter(file, os.Stderr), "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile)
	return file, nil
}
