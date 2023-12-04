package utils

import (
	"errors"
	global "github.com/lowk3v/dumpsc/config"
	"os"
	"strings"
)

func DirExists(dir string, create bool) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if create {
			err := os.MkdirAll(dir, 0755)
			if HandleError(err, "Error create directory") {
				return false
			}
			return true
		}
		global.Log.
			WithField("output", dir).
			Error("Directory is not exist")
		return false
	}
	return true
}

func WriteFile(path string, output string) error {
	if output == "" {
		return errors.New("output path is empty")
	}

	// if path contains dir, mkdir -p
	dir := path[:len(path)-len(path[strings.LastIndex(path, "/"):])]
	if !DirExists(dir, true) {
		return errors.New("error create directory")
	}

	// Write file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(output)
	if err != nil {
		return err
	}

	// Save changes
	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}
