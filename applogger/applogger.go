package applogger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type AppLogger struct {
	Logfile *os.File
	Logger  *log.Logger
}

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

func (l *AppLogger) Close() error {
	return l.Logfile.Close()
}
