/*
Package models defines the models for config file, directories and files, and
defines the methods to create backups and remove the original files.
*/
package models

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aldenso/gotoolbackup/applogger"
)

//Tomlconfig struct to read config file and get parameters
type Tomlconfig struct {
	Title       string
	Directories map[string]Directory
}

//Directory struct indicating origin, destiny directories and a retention time
//in days.
type Directory struct {
	ORIGIN    string
	DESTINY   string
	RETENTION int
}

//Filestobackup struct to check wich files needs backup according to the
//retention time.
type Filestobackup struct {
	ORIGIN  string
	FILES   []string
	DESTINY string
}

//Backups struct to define files associated with origin and destiny directories
//that needs backup.
type Backups struct {
	Elements []Filestobackup
}

//Size method to get the size of files needing backup
func (f *Filestobackup) Size() int64 {
	var size int64
	for _, file := range f.FILES {
		info, _ := os.Stat(f.ORIGIN + "/" + file)
		increment := info.Size()
		size = size + increment
	}
	return size
}

//checkError function to help with error validation and logs
func checkError(l *applogger.AppLogger, err error) {
	if err != nil {
		l.Logger.Println("Error:", err)
		fmt.Println("Error:", err)
		l.Close()
		os.Exit(1)
	}
}

//printLog function to help with print statements and logs
func printLog(l *applogger.AppLogger, msg string) {
	l.Logger.Println(msg)
	fmt.Println(msg)
}

//BackingUP method to create backups with tar and gzip
func (b *Backups) BackingUP(l *applogger.AppLogger) {
	nowref := time.Now()
	backupfilename := string("/backup_" + strings.Replace(nowref.Format(time.RFC3339), ":", "", -1) + ".tar.gz")
	var wg sync.WaitGroup
	for _, v := range b.Elements {
		wg.Add(1)
		go func(v Filestobackup) {
			defer wg.Done()
			backupfile, err := os.Create(v.DESTINY + backupfilename)
			if err != nil {
				checkError(l, err)
			}
			defer backupfile.Close()
			gw := gzip.NewWriter(backupfile)
			defer gw.Close()
			tw := tar.NewWriter(gw)
			defer tw.Close()
			for _, file := range v.FILES {
				openfile, err := os.Open(v.ORIGIN + "/" + file)
				if err != nil {
					checkError(l, err)
				}
				defer openfile.Close()
				if stat, err := openfile.Stat(); err == nil {
					header, err := tar.FileInfoHeader(stat, stat.Name())
					if err != nil {
						checkError(l, err)
					}
					if err := tw.WriteHeader(header); err != nil {
						checkError(l, err)
					}
					if _, err := io.Copy(tw, openfile); err != nil {
						checkError(l, err)
					}
				}
			}
			backupfileSize, _ := os.Stat(v.DESTINY + backupfilename)
			msg := "backup file: " + v.DESTINY + backupfilename + " - size in bytes: " + strconv.FormatInt(backupfileSize.Size(), 10)
			printLog(l, msg)
		}(v)
	}
	wg.Wait()
}

// RemoveOriginalFiles method to delete original files if keepfiles in main is false, only after
// the backup is completed without errors.
func (b *Backups) RemoveOriginalFiles() error {
	for _, v := range b.Elements {
		if len(v.FILES) > 0 {
			for _, file := range v.FILES {
				err := os.Remove(v.ORIGIN + "/" + file)
				if err != nil {
					fmt.Println("failed to remove old files.")
					return err
				}
			}
			fmt.Printf("Removed Original Files for %s: %v\n", v.ORIGIN, v.FILES)
		}
	}
	return nil
}
