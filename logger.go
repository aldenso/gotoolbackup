package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/afero"
)

//AppLogger struct to use a logfile and a logger
type AppLogger struct {
	Logfile afero.File
	Logger  *log.Logger
}

//NewLogger create the new logger
func NewLogger(fs afero.Fs, filename string) *AppLogger {
	file, err := fs.Create(filename + "_" + strings.Replace(time.Now().Format(time.RFC3339), ":", "", -1) + ".log")
	if err != nil {
		fmt.Println("Error creating logger file", err)
		os.Exit(1)
	}
	logger := log.New(file, "gotoolbackup: ", log.Ldate|log.LstdFlags)
	newlogger := &AppLogger{file, logger}
	return newlogger
}

//Close method to close logfile when it finish or fail
func (l *AppLogger) Close() error {
	err := l.Logfile.Close()
	return err
}

// checkError function to help with error validation and logs
func checkError(err error) {
	if err != nil {
		Logs.Logger.Println("Error:", err)
		fmt.Println("Error:", err)
		Logs.Close()
		os.Exit(1)
	}
}

// printLog function to help with print statements and logs
func printLog(msg string) {
	Logs.Logger.Println(msg)
	fmt.Println(msg)
}
