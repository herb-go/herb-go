package project

import (
	"errors"
	"path/filepath"

	"github.com/herb-go/util/cli/app/tools"
)

//ErrNotInHerbGoAPPFolder err rasied when not in herb go app folder
var ErrNotInHerbGoAPPFolder = errors.New("current folder is not a herb-go app folder")

//GetModuleFolder get module folder with given.
//Return module folder and any error rasied.
var GetModuleFolder = func(path string) (string, error) {
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

	var ModulesFolders = []string{
		filepath.Join("src", "modules"),
		filepath.Join("src", "vendor", "modules"),
	}
	for _, v := range FilesMustExists {
		result, err := tools.FileExists(v...)
		if err != nil {
			return "", err
		}
		if result == false {
			return "", ErrNotInHerbGoAPPFolder
		}
	}
	for _, v := range FoldersMustExists {
		result, err := tools.IsFolder(v...)
		if err != nil {
			return "", err
		}
		if result == false {
			return "", ErrNotInHerbGoAPPFolder
		}

	}
	for _, v := range ModulesFolders {
		result, err := tools.IsFolder(path, v)

		if err != nil {
			return "", err
		}
		if result == true {
			return v, nil
		}
	}
	return "", ErrNotInHerbGoAPPFolder
}
