package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aldenso/gotoolbackup/backupfs"
	"github.com/spf13/afero"
)

var config *Tomlconfig

//var tomlfile = "gobackup.toml"
var files = []string{"fileX.txt", "fileY.txt", "fileZ.txt"}
var NewFs = backupfs.InitMemFs()

func readConfig() *Tomlconfig {
	config, err := readTomlFile(*tomlfile)
	if err != nil {
		fmt.Printf("failed to read tomlfile %s\n", err)
		os.Exit(1)
	}
	return config
}

func createMockData() {
	config = readConfig()

	for _, v := range config.Directories {
		err := NewFs.MkdirAll(v.ORIGIN, 0755)
		if err != nil {
			fmt.Println("can't create mock dirs")
		}
		for _, f := range files {
			filepath := v.ORIGIN + "/" + f
			newfile, err := NewFs.Create(filepath)
			if err != nil {
				fmt.Println("can't create mock files")
			}
			_, err = newfile.WriteString("some content")
			if err != nil {
				fmt.Println("can't write to mock files")
			}
			newfile.Close()
			const shortForm = "2006-Jan-02"
			timestring, _ := time.Parse(shortForm, "2013-Feb-03")
			errtime := afero.Fs(NewFs).Chtimes(filepath, timestring, timestring)
			if errtime != nil {
				fmt.Println("can't change time to mock files")
			}
		}
		err = Fs.MkdirAll(v.DESTINY, 0755)
		if err != nil {
			fmt.Println("can't create mock dirs")
		}
	}
}

func Test_BackingUp(t *testing.T) {
	createMockData()
	backup := &Backups{}
	for _, directory := range config.Directories {
		element := checkFiles(NewFs, directory.ORIGIN, directory.DESTINY, directory.RETENTION)
		fmt.Printf("%s\n%s\n", element.ORIGIN, element.FILES)
		if len(element.FILES) == 0 {
			printLog("nothing to backup in: " + element.ORIGIN)
		} else {
			backup.Elements = append(backup.Elements, *element)
		}
		fmt.Println("====================================================")
	}
	if len(backup.Elements) == 0 {
		t.Errorf("Error checking expiration")
	}
	err := backup.CheckFilesPerms(NewFs)
	if err != nil {
		t.Errorf("error checking files perms %s", err)
	}
	for _, i := range backup.Elements {
		files := strings.Join(i.FILES, ",")
		fmt.Println(i.ORIGIN + ": " + files + " - size in bytes: " +
			strconv.FormatInt(i.Size(NewFs), 10))
	}
	errs := backup.BackingUP(NewFs)
	if len(errs) != 0 {
		t.Errorf("Backup ended with errors")
	}
}
