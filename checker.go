package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

//LineSeparator just for separate output in more readable lines
func LineSeparator() {
	fmt.Println("#####################################")
}

//ReadTomlFile function to read the config file
func ReadTomlFile(tomlfile string) (*Tomlconfig, error) {
	var config *Tomlconfig
	if _, err := toml.DecodeFile(tomlfile, &config); err != nil {
		return nil, err
	}
	return config, nil
}

// function to check if origin and destiny exist
func checkExist(origin string, destiny string) {
	LineSeparator()
	dirs := []string{origin, destiny}
	fmt.Println("Checking directories:")
	for _, d := range dirs {
		// if _, err := Fs.Stat(d); os.IsNotExist(err) {
		_, err := Fs.Stat(d)
		if err != nil {
			fmt.Printf("FAILED: %s\n%v\n", d, err.Error())
			os.Exit(1)
			// return
		}
		fmt.Printf("PASS: %s\n", d)
	}
	// return
}

//RunCheck function to run check
func RunCheck(config Tomlconfig) {
	LineSeparator()
	fmt.Printf("Config Title:\n%s\n", config.Title)
	LineSeparator()
	for directoryName, directory := range config.Directories {
		fmt.Printf("Backup: %s\nOrigin: %s | Destiny: %s | Retention: %d\n", directoryName, directory.ORIGIN, directory.DESTINY, directory.RETENTION)
	}
	for _, d := range config.Directories {
		checkExist(d.ORIGIN, d.DESTINY)
	}
}

//CheckFiles check wich files needs backup according to retention
func CheckFiles(dirorigin string, dirdestiny string, retention int) *Filestobackup {
	retentionhours := (retention * 24)
	backup := &Filestobackup{}
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
