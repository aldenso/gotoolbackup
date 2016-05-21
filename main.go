/*
Program to create backups using toml file, where you indicate origin and destiny
directories and retention.

TODO: include log package to create output file
*/
package main

import (
	"flag"
	"fmt"
	"gotoolbackup/checkers"
	"gotoolbackup/models"
	"os"
)

var (
	tomlfile = flag.String("tomlfile", "gobackup.toml",
		"indicate tomlfile to read backups details.")
	keepfiles = flag.Bool("keepfiles", true,
		"indicate if you want to keep original files.")
)

func printUsedValues() {
	fmt.Println("#### Running with values ####")
	fmt.Println("tomlfile:", *tomlfile)
	fmt.Println("keepfiles:", *keepfiles)
}

func main() {
	flag.Parse()
	printUsedValues()
	config := checkers.ReadTomlFile(*tomlfile)
	continuebackup := checkers.RunCheck(config)
	if continuebackup != true {
		os.Exit(1)
	}
	checkers.LineSeparator()
	fmt.Printf("%s", "Checking Retention for files\n")
	checkers.LineSeparator()
	backup := &models.Backups{}
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
		if !*keepfiles {
			err := backup.RemoveOriginalFiles()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
