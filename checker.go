package main

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/spf13/afero"
)

// LineSeparator just for separate output in more readable lines
func LineSeparator() {
	fmt.Println("#####################################")
}

// ReadTomlFile function to read the config file
func readTomlFile(tomlfile string) (*Tomlconfig, error) {
	var config *Tomlconfig
	if _, err := toml.DecodeFile(tomlfile, &config); err != nil {
		return nil, err
	}
	return config, nil
}

// function to check if origin and destiny exist
func checkExist(fs afero.Fs, directory string) error {
	LineSeparator()
	_, err := fs.Stat(directory)
	if err != nil {
		return err
	}
	return nil
}

// runCheck function to run check
func runCheck(fs afero.Fs, config Tomlconfig) {
	LineSeparator()
	fmt.Printf("Config Title:\n%s\n", config.Title)
	LineSeparator()
	for directoryName, directory := range config.Directories {
		fmt.Printf("Backup: %s\nOrigin: %s | Destiny: %s | Retention: %d\n", directoryName, directory.ORIGIN, directory.DESTINY, directory.RETENTION)
	}
	fmt.Println("Checking directories:")
	var numerrs int
	for _, d := range config.Directories {
		err := checkExist(fs, d.ORIGIN)
		if err != nil {
			numerrs++
			checkError(err)
		}
		fmt.Printf("PASS: %s\n", d.ORIGIN)
		err = checkExist(fs, d.DESTINY)
		if err != nil {
			numerrs++
			checkError(err)
		}
		fmt.Printf("PASS: %s\n", d.DESTINY)
	}
	if numerrs != 0 {
		printLog("Please create all directories")
		os.Exit(1)
	}
}

// checkFiles check wich files needs backup according to retention
func checkFiles(fs afero.Fs, dirorigin string, dirdestiny string, retention int) *Filestobackup {
	retentionhours := (retention * 24)
	backup := &Filestobackup{}
	backup.ORIGIN = dirorigin
	backup.DESTINY = dirdestiny
	files, err := afero.ReadDir(fs, dirorigin)
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
