package applogger

import (
	"fmt"
	"log"
	"os"
	"time"
)

//AppLogger struct to use a logfile and a logger
type AppLogger struct {
	Logfile *os.File
	Logger  *log.Logger
}

//NewLogger create the new logger
func NewLogger(filename string) *AppLogger {
	file, err := os.Create(filename + "_" + time.Now().Format(time.RFC3339) + ".log")
	if err != nil {
		fmt.Println("Error creating logger file", err)
		os.Exit(1)
	}
	logger := log.New(file, "gotoolbackup: ", log.Ldate|log.LstdFlags)
	newlogger := &AppLogger{file, logger}
	return newlogger
}

//Close function to close logfile when it finish or fail
func (l *AppLogger) Close() error {
	err := l.Logfile.Close()
	return err
}
