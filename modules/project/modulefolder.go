package project

import (
	"strings"
	"path/filepath"
	"errors"
	"io/ioutil"
	"github.com/herb-go/util/cli/app/tools"
)

var ErrNotInHerbGoAPPFolder = errors.New("current folder is not a herb-go app folder")

var GetModuleFolder = func(path string) (string,error) {
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

		var ModulesFolders=[]string{
		"src/modules",
		"src/vendor/modules",
	}
	for _, v := range FilesMustExists {
		result, err := tools.FileExists(v...)
		if err != nil {
			return "",err
		}
		if result == false {
			return "",ErrNotInHerbGoAPPFolder
		}
	}
	for _, v := range FoldersMustExists {
		result, err := tools.IsFolder(v...)
		if err != nil {
			return "",err
		}
		if result == false {
			return "",ErrNotInHerbGoAPPFolder
		}

	}
	for _,v:=range ModulesFolders{
		gomodpath:=filepath.Join(path,v,"go.mod")
		result,err:=tools.FileExists(gomodpath)
		if err!=nil{
			return "",err
		}
		if result==true {
			bs,err:=ioutil.ReadFile(gomodpath)
			if err!=nil{
				return "",err
			}
		lines:=strings.Split(string(bs),"\n")
		if len(lines)==0{
			continue
		}
		words:=strings.Split(strings.TrimSpace(lines[0])," ")
		if len(words) <2||words[0]!="module" || words[1]!="modules"{
			continue
		}
		return v,nil
		}
	}
	return "",ErrNotInHerbGoAPPFolder
}
