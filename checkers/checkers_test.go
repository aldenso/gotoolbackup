package checkers

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aldenso/gotoolbackup/models"
)

func TestRunCheck(t *testing.T) {
	dirs1 := &models.Directory{
		ORIGIN:    "/",
		DESTINY:   "/var",
		RETENTION: 15,
	}
	dirs2 := &models.Directory{
		ORIGIN:    "/home",
		DESTINY:   "/tmp",
		RETENTION: 30,
	}
	tomlconfig := &models.Tomlconfig{
		Title:       "test",
		Directories: map[string]models.Directory{"DIRNAME1": *dirs1, "DIRNAME2": *dirs2},
	}

	if err := RunCheck(*tomlconfig); err != nil {
		t.Errorf("Failed to in RunCheck, got error '%v'", err)
	}
}

func TestCheckFiles(t *testing.T) {
	if err := os.Mkdir("/tmp/gotoolbackup", 0750); err != nil {
		fmt.Println("Can't create temp dir", err)
	}
	testfileA, err := os.Create("/tmp/gotoolbackup/testfileA.txt")
	if err != nil {
		fmt.Println("Error creating log file", err)
	}
	testfileB, err := os.Create("/tmp/gotoolbackup/testfileB.txt")
	if err != nil {
		fmt.Println("Error creating log file", err)
	}
	err = os.Chtimes("/tmp/gotoolbackup/testfileB.txt", time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, -1, 0))
	if err != nil {
		fmt.Println("Can't change timestamp to /tmp/gotoolbackup/testfileB.txt")
	}
	defer testfileA.Close()
	defer testfileB.Close()
	dirs := &models.Directory{
		ORIGIN:    "/tmp/gotoolbackup",
		DESTINY:   "/var",
		RETENTION: 15,
	}
	filestobackup := CheckFiles(dirs.ORIGIN, dirs.DESTINY, dirs.RETENTION)
	if err = os.Remove("/tmp/gotoolbackup/testfileA.txt"); err != nil {
		fmt.Println("Can't remove testfile", err)
	}
	if err = os.Remove("/tmp/gotoolbackup/testfileB.txt"); err != nil {
		fmt.Println("Can't remove testfile", err)
	}
	if err = os.RemoveAll("/tmp/gotoolbackup"); err != nil {
		fmt.Println("Can't remove dir /tmp/gotoolbackup", err)
	}
	if filestobackup.FILES[0] != "testfileB.txt" {
		t.Errorf("Error Checking files, expected 'testfileB.txt', got '%s'", filestobackup.FILES[0])
	}
}
