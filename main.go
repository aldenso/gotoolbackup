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

	"github.com/aldenso/gotoolbackup/backupfs"
	"github.com/spf13/afero"
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

// Fs afero fs to help later with testing.
// var Fs = afero.NewOsFs()
var Fs = backupfs.InitOSFs()

// NowRef to use as pattern for backup file names.
var NowRef = time.Now()

// function to indicate the values you are using
func printUsedValues() {
	fmt.Println("#### Running with values ####")
	fmt.Println("tomlfile:", *tomlfile)
	fmt.Println("remove:", *removefiles)
	fmt.Println("log:", *logfile)
}

// Logs logger defined in logger
var Logs *AppLogger

// RunBackups function to run all backups from main
func RunBackups(fs afero.Fs) []error {
	var errs []error
	start := time.Now()
	flag.Parse()
	printUsedValues()
	Logs = NewLogger(fs, *logfile)
	printLog("Reading tomlfile: " + *tomlfile)
	config, err := readTomlFile(*tomlfile)
	checkError(err)
	runCheck(*config)
	LineSeparator()
	printLog("Checking Retention for files")
	LineSeparator()
	backup := &Backups{}
	for _, directory := range config.Directories {
		element := checkFiles(fs, directory.ORIGIN, directory.DESTINY, directory.RETENTION)
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
	err = backup.CheckFilesPerms(fs)
	checkError(err)
	printLog("Running backups for: ")
	for _, i := range backup.Elements {
		files := strings.Join(i.FILES, ",")
		printLog(i.ORIGIN + ": " + files + " - size in bytes: " +
			strconv.FormatInt(i.Size(fs), 10))
	}
	msgs, errs := backup.BackingUP(fs)
	if len(errs) == 0 {
		printLog("Backup Successful")
	} else {
		printLog("Backup Ended with errors: ")
		for _, e := range errs {
			printLog(e.Error())
			errs = append(errs, e)
		}
	}
	for _, msg := range msgs {
		printLog(msg)
	}
	if *removefiles {
		filelist, delerr := backup.RemoveOriginalFiles(fs)
		if delerr != nil {
			fmt.Println("failed to remove some of the old files.")
			for _, file := range filelist {
				fmt.Printf("failed to remove: %s\n", file)
				errs = append(errs, delerr)
			}
		}
		printLog("removing old files ended!")
	}
	printLog("gotoolbackup ended! in: " + time.Since(start).String())
	err = Logs.Close()
	if err != nil {
		fmt.Printf("Error closing logger: %v", err)
		errs = append(errs, err)
	}
	return errs
}

func main() {
	errs := RunBackups(Fs)
	if len(errs) != 0 {
		fmt.Printf("Errors found in backups\n")
		for _, err := range errs {
			fmt.Printf("%s\n", err)
		}
	}
}
