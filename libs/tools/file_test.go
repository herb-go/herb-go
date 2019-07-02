package tools

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestFile(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "herb-go-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(tmpdir)
	}()
	result, err := FileExists(path.Join(tmpdir, "notexists"))
	if err != nil {
		t.Fatal(err)
	}
	if result {
		t.Fatal(result)
	}
	ioutil.WriteFile(path.Join(tmpdir, "file"), []byte("123"), 0700)
	result, err = FileExists(path.Join(tmpdir, "file"))
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Fatal(result)
	}
}
