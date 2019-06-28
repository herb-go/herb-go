package tools

import (
	"io/ioutil"
	"os"
	"strings"
)

func ReplaceLine(path string, from string, to string) (bool, error) {
	var found bool
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	lines := strings.Split(string(data), "\r\n")
	var line string
	for k := range lines {
		line = lines[k]
		if strings.TrimSpace(line) == from {
			found = true
			if line != "" {
				lines[k] = to
			}
		}

	}
	if found {
		output := strings.Join(lines, "\r\n")
		err := ioutil.WriteFile(path, []byte(output), info.Mode())
		if err != nil {
			return false, (err)
		}
		return true, nil
	}
	return false, nil
}
