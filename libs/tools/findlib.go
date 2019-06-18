package tools

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

var ErrGoEnvEmpty = errors.New("Go path env is empty.Please set GOPATH env in yonr shell config.")

var NewLibNotFoundError = func(lib string) error {
	return fmt.Errorf(" \"%s\" not found.\nYou should use \"go get -u %s\" to install", lib, lib)
}

func FindLib(gopathenv string, libname string) (string, error) {
	var result string
	if gopathenv == "" {
		return "", ErrGoEnvEmpty
	}
	pathlist := strings.Split(gopathenv, ":")
	for _, v := range pathlist {
		result = path.Join(v, "src", libname)
		_, err := os.Stat(result)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}
		return result, nil
	}
	return "", NewLibNotFoundError(libname)
}
