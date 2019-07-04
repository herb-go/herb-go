package config

import (
	"errors"

	"github.com/herb-go/util/cli/app/tools"
)

var ErrNotInHerbGoAPPFolder = errors.New("current folder is not a herb-go app folder")

func ErrorIfNotInAppFolder(path string) error {
	var FilesMustExists = [][]string{
		[]string{path, "src", "main.go"},
	}
	var FoldersMustExists = [][]string{
		[]string{path, "config"},
		[]string{path, "appdata"},
		[]string{path, "resources"},
		[]string{path, "system", "config.examples"},
		[]string{path, "system", "constants"},
	}
	for _, v := range FilesMustExists {
		result, err := tools.FileExists(v...)
		if err != nil {
			return err
		}
		if result == false {
			return ErrNotInHerbGoAPPFolder
		}
	}
	for _, v := range FoldersMustExists {
		result, err := tools.IsFolder(v...)
		if err != nil {
			return err
		}
		if result == false {
			return ErrNotInHerbGoAPPFolder
		}

	}
	return nil
}
