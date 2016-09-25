/*
Program to create backups using toml file, where you indicate origin and destiny
directories and retention.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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

//Logs logger defined in logger package
var Logs *applogger.AppLogger

//checkError function to help with error validation and logs
func checkError(err error) {
	if err != nil {
		Logs.Logger.Println("Error:", err)
		fmt.Println("Error:", err)
		Logs.Close()
		os.Exit(1)
	}
}

//printLog function to help with print statements and logs
func printLog(msg string) {
	Logs.Logger.Println(msg)
	fmt.Println(msg)
}

func main() {
	start := time.Now()
	flag.Parse()
	printUsedValues()
	Logs = applogger.NewLogger(*logfile)
	printLog("Reading tomlfile: " + *tomlfile)
	config, err := checkers.ReadTomlFile(*tomlfile)
	checkError(err)
	err = checkers.RunCheck(*config)
	if err != nil {
		checkError(err)
	}
	checkers.LineSeparator()
	printLog("Checking Retention for files")
	checkers.LineSeparator()
	backup := &models.Backups{}
	for _, directory := range config.Directories {
		element := checkers.CheckFiles(directory.ORIGIN, directory.DESTINY, directory.RETENTION)
		fmt.Printf("%s\n%s\n", element.ORIGIN, element.FILES)
		if len(element.FILES) == 0 {
			printLog("nothing to backup in: " + element.ORIGIN)
		} else {
			backup.Elements = append(backup.Elements, *element)
		}
		fmt.Println("====================================================")
	}
	if len(backup.Elements) == 0 {
		printLog("None of the files needs a backup!")
		Logs.Close()
		os.Exit(0)
	}
	printLog("Running backups for: ")
	for _, i := range backup.Elements {
		files := strings.Join(i.FILES, ",")
		printLog(i.ORIGIN + ": " + files + " - size in bytes: " +
			strconv.FormatInt(i.Size(), 10))
	}
	backup.BackingUP(Logs)
	printLog("Backup Successful")
	if *removefiles {
		err = backup.RemoveOriginalFiles()
		checkError(err)
		printLog("old files removed")
	}
	printLog("gotoolbackup ended! in: " + time.Since(start).String())
	err = Logs.Close()
	checkError(err)
}
