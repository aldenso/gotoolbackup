/*
Program to create backups using toml file, where you indicate origin and destiny
directories and retention.

TODO: use external tomlfile, use flags to indicate if you want to erase original
files or not.
*/
package main

import (
	"fmt"
	"gotoolbackup/checkers"
	"gotoolbackup/models"
	"os"

	"github.com/BurntSushi/toml"
)

func main() {
	continuebackup := checkers.RunCheck()
	if continuebackup != true {
		os.Exit(1)
	}
	checkers.LineSeparator()
	fmt.Printf("%s", "Checking Retention for files\n")
	checkers.LineSeparator()
	backup := &models.Backups{}
	//  needs improvement
	var config models.Tomlconfig
	if _, err := toml.DecodeFile("gobackup.toml", &config); err != nil {
		fmt.Println(err)
	}
	// needs improvement
	for _, directory := range config.Directories {
		salida := checkers.CheckFiles(directory.ORIGIN, directory.DESTINY, directory.RETENTION)
		fmt.Printf("%s\n%s\n", salida.ORIGIN, salida.FILES)
		if len(salida.FILES) == 0 {
			fmt.Println("nothing to backup in:", salida.ORIGIN)
		} else {
			backup.Elements = append(backup.Elements, *salida)
		}
		fmt.Println("====================================================")
	}
	fmt.Println("BACKING", backup.Elements)
	err := backup.BackingUP()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Backup Successful")
	}
}
