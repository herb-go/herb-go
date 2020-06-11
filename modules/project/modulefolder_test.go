package project

import (
	"path/filepath"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/herb-go/util"
)

func TestGetModuleFolder(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "herb-go-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(tmpdir)
	}()
	mpath,err :=GetModuleFolder(tmpdir)
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
		[]string{tmpdir, "system", "configskeleton"},
		[]string{tmpdir, "system", "constants"},
	}
	for _, v := range FoldersMustExists {
		mpath,err = GetModuleFolder(tmpdir)
		if err != ErrNotInHerbGoAPPFolder {
			t.Fatal(err)
		}
		err = os.MkdirAll(path.Join(v...), util.DefaultFolderMode)
		if err != nil {
			t.Fatal(err)
		} 
	}
	mpath,err = GetModuleFolder(tmpdir)
	if err != ErrNotInHerbGoAPPFolder {
		t.Fatal(err)
	}
	gomodpath:=filepath.Join(tmpdir,"src","modules","go.mod")
	err=os.MkdirAll(filepath.Dir(gomodpath),util.DefaultFolderMode)
	if err != nil {
		t.Fatal(err)
	}
	err=ioutil.WriteFile(gomodpath,[]byte{},util.DefaultFileMode)
	if err != nil {
		t.Fatal(err)
	}
	mpath,err = GetModuleFolder(tmpdir)
	if err != ErrNotInHerbGoAPPFolder {
		t.Fatal(err)
	}
	err=ioutil.WriteFile(gomodpath,[]byte("module modules"),util.DefaultFileMode)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	mpath,err = GetModuleFolder(tmpdir)
	if err != nil {
		t.Fatal(err)
	}
	if mpath!="src/modules"{
		t.Fatal(mpath)
	}
	err=os.Remove(gomodpath)
	if err!=nil{
		t.Fatal(err)
	}
		mpath,err = GetModuleFolder(tmpdir)
	if err != ErrNotInHerbGoAPPFolder {
		t.Fatal(err)
	}
	gomodpath=filepath.Join(tmpdir,"src","vendor","modules","go.mod")
	err=os.MkdirAll(filepath.Dir(gomodpath),util.DefaultFolderMode)
	if err != nil {
		t.Fatal(err)
	}
	err=ioutil.WriteFile(gomodpath,[]byte{},util.DefaultFileMode)
	if err != nil {
		t.Fatal(err)
	}
	mpath,err = GetModuleFolder(tmpdir)
	if err != ErrNotInHerbGoAPPFolder {
		t.Fatal(err)
	}
	err=ioutil.WriteFile(gomodpath,[]byte("module modules"),util.DefaultFileMode)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	mpath,err = GetModuleFolder(tmpdir)
	if err != nil {
		t.Fatal(err)
	}
	if mpath!="src/vendor/modules"{
		t.Fatal(mpath)
	}
}
