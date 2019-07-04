package config

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/herb-go/util"
)

func TestCheckApp(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "herb-go-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(tmpdir)
	}()
	err = ErrorIfNotInAppFolder(tmpdir)
	if err != ErrNotInHerbGoAPPFolder {
		t.Fatal(err)
	}
	err = os.MkdirAll(path.Join(tmpdir, "src"), util.DefaultFolderMode)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(path.Join(tmpdir, "src", "main.go"), []byte("123"), util.DefaultFileMode)
	if err != nil {
		t.Fatal(err)
	}
	var FoldersMustExists = [][]string{
		[]string{tmpdir, "config"},
		[]string{tmpdir, "appdata"},
		[]string{tmpdir, "resources"},
		[]string{tmpdir, "system", "config.examples"},
		[]string{tmpdir, "system", "constants"},
	}
	for _, v := range FoldersMustExists {
		err = ErrorIfNotInAppFolder(tmpdir)
		if err != ErrNotInHerbGoAPPFolder {
			t.Fatal(err)
		}
		err = os.MkdirAll(path.Join(v...), util.DefaultFolderMode)
		if err != nil {
			t.Fatal(err)
		}
	}
	err = ErrorIfNotInAppFolder(tmpdir)
	if err != nil {
		t.Fatal(err)
	}
}
