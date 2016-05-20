package models

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"time"
)

// struct to read file to get parameters
type Tomlconfig struct {
	Title       string
	Directories map[string]Directory
}

// parameters taken from file
type Directory struct {
	ORIGIN    string
	DESTINY   string
	RETENTION int
}

// struct to check wich files need backup
type Filestobackup struct {
	ORIGIN  string
	FILES   []string
	DESTINY string
}

// files associated with origin and destiny that needs backup
type Backups struct {
	Elements []Filestobackup
}

// method to create backups with tar and gzip
func (b *Backups) BackingUP() error {
	for _, v := range b.Elements {
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
