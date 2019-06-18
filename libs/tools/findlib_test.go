package tools

import (
	"os"
	"strings"
	"testing"
)

func TestFindLib(t *testing.T) {
	var result string
	var err error
	result, err = FindLib("", "github.com/herb-go/herb-go/libs/tools")
	if result != "" {
		t.Fatal(result)
	}
	if err != ErrGoEnvEmpty {
		t.Fatal(err)
	}
	result, err = FindLib(strings.Join([]string{os.Getenv("GOPATH"), ""}, ":"), "github.com/herb-go/herb-go/notexists/libs/tools")
	if result != "" {
		t.Fatal(result)
	}
	if err == nil {
		t.Fatal(err)
	}
	result, err = FindLib(strings.Join([]string{os.Getenv("GOPATH"), ""}, ":"), "github.com/herb-go/herb-go/libs/tools")
	if result == "" {
		t.Fatal(result)
	}
	if err != nil {
		t.Fatal(err)
	}
}
