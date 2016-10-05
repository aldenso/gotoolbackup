package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/afero"
)

// Tomlconfig struct to read config file and get parameters
type Tomlconfig struct {
	Title       string
	Directories map[string]Directory
}

// Directory struct indicating origin, destiny directories and a retention time
// in days.
type Directory struct {
	ORIGIN    string
	DESTINY   string
	RETENTION int
}

// Filestobackup struct to check wich files needs backup according to the
// retention time.
type Filestobackup struct {
	ORIGIN  string
	FILES   []string
	DESTINY string
}

// Backups struct to define files associated with origin and destiny directories
// that needs backup.
type Backups struct {
	Elements []Filestobackup
}

// Size method to get the size of files needing backup
func (f *Filestobackup) Size(fs afero.Fs) int64 {
	var size int64
	for _, file := range f.FILES {
		info, _ := fs.Stat(f.ORIGIN + "/" + file)
		increment := info.Size()
		size = size + increment
	}
	return size
}

// BackingUP method to create backups with tar and gzip
func (b *Backups) BackingUP(fs afero.Fs) ([]string, []error) {
	var msgs []string
	var errs []error
	backupfilename := string("/backup_" + strings.Replace(NowRef.Format(time.RFC3339), ":", "", -1) + ".tar.gz")
	var wg sync.WaitGroup
	for _, v := range b.Elements {
		wg.Add(1)
		go func(v Filestobackup) {
			defer wg.Done()
			backupfile, err := fs.Create(v.DESTINY + backupfilename)
			if err != nil {
				errs = append(errs, err)
			}
			defer backupfile.Close()
			gw := gzip.NewWriter(backupfile)
			defer gw.Close()
			tw := tar.NewWriter(gw)
			defer tw.Close()
			for _, file := range v.FILES {
				openfile, err := fs.Open(v.ORIGIN + "/" + file)
				if err != nil {
					errs = append(errs, err)
				}
				if stat, err := openfile.Stat(); err == nil {
					header, err := tar.FileInfoHeader(stat, stat.Name())
					if err != nil {
						errs = append(errs, err)
					}
					if err := tw.WriteHeader(header); err != nil {
						checkError(err)
						//errs = append(errs, err)
					}
					if _, err := io.Copy(tw, openfile); err != nil {
						errs = append(errs, err)
					}
				}
				openfile.Close()
			}
			backupfileSize, _ := fs.Stat(v.DESTINY + backupfilename)
			msg := "backup file: " + v.DESTINY + backupfilename + " - size in bytes: " + strconv.FormatInt(backupfileSize.Size(), 10)
			msgs = append(msgs, msg)
		}(v)
	}
	wg.Wait()
	return msgs, errs
}

// CheckFilesPerms checks files to backup perms before backup
func (b *Backups) CheckFilesPerms(fs afero.Fs) error {
	for _, v := range b.Elements {
		for _, file := range v.FILES {
			_, err := fs.Open(v.ORIGIN + "/" + file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// RemoveOriginalFiles method to delete original files if keepfiles in main is false, only after
// the backup is completed without errors.
func (b *Backups) RemoveOriginalFiles() ([]string, error) {
	var missing []string
	var reterr error
	for _, v := range b.Elements {
		fmt.Printf("Removing Original Files for %s: %v\n", v.ORIGIN, v.FILES)
		if len(v.FILES) > 0 {
			for _, file := range v.FILES {
				err := Fs.Remove(v.ORIGIN + "/" + file)
				if err != nil {
					fullfile := fmt.Sprintf("%s/%s", v.ORIGIN, file)
					missing = append(missing, fullfile)
					reterr = err
				}
			}
		}
	}
	if len(missing) != 0 {
		return missing, reterr
	}
	return nil, nil
}
