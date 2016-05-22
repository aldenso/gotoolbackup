package checkers

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gotoolbackup/models"

	"github.com/BurntSushi/toml"
)

func LineSeparator() {
	fmt.Println("#####################################")
}

func ReadTomlFile(tomlfile string) (*models.Tomlconfig, error) {
	var config *models.Tomlconfig
	if _, err := toml.DecodeFile(tomlfile, &config); err != nil {
		return nil, err
	}
	return config, nil
}

// function to check if origin and destiny exist
func checkExist(origin string, destiny string) bool {
	var dirErrors bool
	LineSeparator()
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
func RunCheck(config models.Tomlconfig) bool {
	continuebackup := true
	LineSeparator()
	fmt.Printf("Config Title:\n%s\n", config.Title)
	LineSeparator()
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
func CheckFiles(dirorigin string, dirdestiny string, retention int) *models.Filestobackup {
	retentionhours := (retention * 24)
	backup := &models.Filestobackup{}
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
