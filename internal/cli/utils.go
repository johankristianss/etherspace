package cli

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/johankristianss/evrium/pkg/build"
	log "github.com/sirupsen/logrus"
)

func truncateString(s string) string {
	if len(s) > MAX_VALUE_LENGTH {
		return s[:MAX_VALUE_LENGTH] + "..."
	}
	return s
}

func checkIfDirExists(dirPath string) error {
	fileInfo, err := os.Stat(dirPath)
	if err == nil {
		if fileInfo.IsDir() {
			return errors.New(dirPath + " already exists")
		}
	}
	return nil
}

func checkIfDirIsEmpty(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}

	return errors.New("Current directory is not empty, try create a new direcory and retry")
}

func CheckError(err error) {
	if err != nil {
		log.WithFields(log.Fields{"Error": err, "BuildVersion": build.BuildVersion, "BuildTime": build.BuildTime}).Error(err.Error())
		os.Exit(-1)
	}
}

func formatTimestamp(timestamp string) string {
	return strings.Replace(timestamp, "T", " ", 1)
}
