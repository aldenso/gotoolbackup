package main

import (
	"fmt"
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
func checkExist(origin string, destiny string) ([]string, error) {
	LineSeparator()
	var missing []string
	var reterr error
	dirs := []string{origin, destiny}
	for _, d := range dirs {
		_, err := Fs.Stat(d)
		if err != nil {
			missing = append(missing, d)
			reterr = err
		}
	}
	if len(missing) != 0 {
		return missing, reterr
	}
	return nil, nil
}

// runCheck function to run check
func runCheck(config Tomlconfig) {
	LineSeparator()
	fmt.Printf("Config Title:\n%s\n", config.Title)
	LineSeparator()
	for directoryName, directory := range config.Directories {
		fmt.Printf("Backup: %s\nOrigin: %s | Destiny: %s | Retention: %d\n", directoryName, directory.ORIGIN, directory.DESTINY, directory.RETENTION)
	}
	fmt.Println("Checking directories:")
	for _, d := range config.Directories {
		dirs, err := checkExist(d.ORIGIN, d.DESTINY)
		if err != nil {
			for _, d := range dirs {
				fmt.Printf("FAILED: %s\n%v\n", d, err.Error())
			}
		} else {
			for _, d := range dirs {
				fmt.Printf("PASS: %s\n", d)
			}
		}
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
