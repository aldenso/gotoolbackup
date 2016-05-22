/*
Program to create backups using toml file, where you indicate origin and destiny
directories and retention.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aldenso/gotoolbackup/applogger"
	"github.com/aldenso/gotoolbackup/checkers"
	"github.com/aldenso/gotoolbackup/models"
)

// variables to indicate flags values
var (
	tomlfile = flag.String("tomlfile", "gobackup.toml",
		"indicate tomlfile to read backups details.")
	removefiles = flag.Bool("remove", false,
		"indicate if you want to remove original files after backup.")
	logfile = flag.String("log", "gotoolbackup",
		"indicate the log name pattern.")
)

// function to indicate the values you are using
func printUsedValues() {
	fmt.Println("#### Running with values ####")
	fmt.Println("tomlfile:", *tomlfile)
	fmt.Println("remove:", *removefiles)
	fmt.Println("log:", *logfile)
}

// logger defined in logger package
var Logs *applogger.AppLogger

// function to help with error validation and logs
func CheckError(e error) {
	if e != nil {
		Logs.Logger.Println("Error:", e)
		fmt.Println("Error:", e)
		Logs.Close()
		os.Exit(1)
	}
}

// function to help with print statements and logs
func PrintLog(msg string) {
	Logs.Logger.Println(msg)
	fmt.Println(msg)
}

func main() {
	start := time.Now()
	flag.Parse()
	printUsedValues()
	Logs = applogger.NewLogger(*logfile)
	config, err := checkers.ReadTomlFile(*tomlfile)
	CheckError(err)
	continuebackup := checkers.RunCheck(*config)
	if continuebackup != true {
		os.Exit(1)
	}
	checkers.LineSeparator()
	PrintLog("Checking Retention for files")
	checkers.LineSeparator()
	backup := &models.Backups{}
	for _, directory := range config.Directories {
		salida := checkers.CheckFiles(directory.ORIGIN, directory.DESTINY, directory.RETENTION)
		fmt.Printf("%s\n%s\n", salida.ORIGIN, salida.FILES)
		if len(salida.FILES) == 0 {
			PrintLog("nothing to backup in: " + salida.ORIGIN)
		} else {
			backup.Elements = append(backup.Elements, *salida)
		}
		fmt.Println("====================================================")
	}
	if len(backup.Elements) == 0 {
		PrintLog("None of the files needs a backup!")
		Logs.Close()
		os.Exit(0)
	}
	PrintLog("Running backups for: ")
	for _, i := range backup.Elements {
		files := strings.Join(i.FILES, ",")
		PrintLog(i.ORIGIN + ": " + files)
	}
	err = backup.BackingUP()
	CheckError(err)
	PrintLog("Backup Successful")
	if *removefiles {
		err := backup.RemoveOriginalFiles()
		CheckError(err)
		PrintLog("old files removed")
	}
	PrintLog("gotoolbackup ended! in: " + time.Since(start).String())
	Logs.Close()
}
