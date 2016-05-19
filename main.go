/*
Program to create backups using toml file, where you indicate origin and destiny
directories and retention.

TODO: use external tomlfile, use flags to indicate if you want to erase original
files or not.
*/
package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

func lines() {
	fmt.Println("#####################################")
}

// struct to read file to get parameters
type tomlconfig struct {
	Title       string
	Directories map[string]directory
}

// parameters taken from file
type directory struct {
	ORIGIN    string
	DESTINY   string
	RETENTION int
}

// struct to check wich files need backup
type filestobackup struct {
	ORIGIN  string
	FILES   []string
	DESTINY string
}

// files associated with origin and destiny that needs backup
type Backups struct {
	elements []filestobackup
}

// function to create backups with tar and gzip
func (b *Backups) BackingUP() error {
	for _, v := range b.elements {
		backupfile, err := os.Create(v.DESTINY + "/backup_" + time.Now().Format(time.RFC3339) + ".tar.gz")
		if err != nil {
			return err
		}
		defer backupfile.Close()
		gw := gzip.NewWriter(backupfile)
		defer gw.Close()
		tw := tar.NewWriter(gw)
		defer tw.Close()
		for _, file := range v.FILES {
			openfile, err := os.Open(v.ORIGIN + "/" + file)
			if err != nil {
				return err
			}
			defer openfile.Close()
			if stat, err := openfile.Stat(); err == nil {
				header, err := tar.FileInfoHeader(stat, stat.Name())
				if err != nil {
					return err
				}
				if err := tw.WriteHeader(header); err != nil {
					return err
				}
				if _, err := io.Copy(tw, openfile); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// function to check if origin and destiny exist
func checkExist(origin string, destiny string) bool {
	var dirErrors bool
	lines()
	dirs := []string{origin, destiny}
	fmt.Println("Checking directories:")
	for _, d := range dirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			fmt.Printf("FAILED: %s\n%v\n", d, err.Error())
			dirErrors = true
		} else {
			fmt.Printf("PASS: %s\n", d)
		}
	}
	return dirErrors
}

// function to run check
func RunCheck() bool {
	continuebackup := true
	var config tomlconfig
	if _, err := toml.DecodeFile("gobackup.toml", &config); err != nil {
		fmt.Println(err)
		return false
	}
	lines()
	fmt.Printf("Config Title:\n%s\n", config.Title)
	lines()
	for directoryName, directory := range config.Directories {
		fmt.Printf("Backup: %s\nOrigin: %s | Destiny: %s | Retention: %d\n", directoryName, directory.ORIGIN, directory.DESTINY, directory.RETENTION)
	}
	for _, d := range config.Directories {
		filesOK := checkExist(d.ORIGIN, d.DESTINY)
		if filesOK != true {
			fmt.Println("++++ PASS!!!!!!! ++++")
		} else {
			fmt.Println("---- FAILED!!!!! ----")
			continuebackup = false
		}
	}
	return continuebackup
}

// check wich files needs backup according to retention
func CheckFiles(dirorigin string, dirdestiny string, retention int) *filestobackup {
	retentionhours := (retention * 24)
	backup := &filestobackup{}
	backup.ORIGIN = dirorigin
	backup.DESTINY = dirdestiny
	files, err := ioutil.ReadDir(dirorigin)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	for _, file := range files {
		expiration := time.Since(file.ModTime()).Hours()
		if expiration >= float64(retentionhours) {
			backup.FILES = append(backup.FILES, file.Name())
		}
	}
	return backup
}

func main() {
	continuebackup := RunCheck()
	if continuebackup != true {
		os.Exit(1)
	}
	lines()
	fmt.Printf("%s", "Checking Retention for files\n")
	lines()
	backup := &Backups{}
	//  needs improvement
	var config tomlconfig
	if _, err := toml.DecodeFile("gobackup.toml", &config); err != nil {
		fmt.Println(err)
	}
	// needs improvement
	for _, directory := range config.Directories {
		salida := CheckFiles(directory.ORIGIN, directory.DESTINY, directory.RETENTION)
		fmt.Printf("%s\n%s\n", salida.ORIGIN, salida.FILES)
		if len(salida.FILES) == 0 {
			fmt.Println("nothing to backup in:", salida.ORIGIN)
		} else {
			backup.elements = append(backup.elements, *salida)
		}
		fmt.Println("====================================================")
	}
	fmt.Println("BACKING", backup.elements)
	err := backup.BackingUP()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Backup Successful")
	}
}
