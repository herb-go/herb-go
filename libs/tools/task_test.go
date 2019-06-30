package tools

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestTask(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "herb-go-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(tmpdir)
	}()
	renderdata := map[string]interface{}{
		"data": "data",
	}
	task := NewTask(path.Join("./", "testdata"), tmpdir)
	err = task.Copy("/demo.txt", "/demo.txt")
	if err != nil {
		t.Fatal(err)
	}
	err = task.Render("/demo1.tmpl", "demo1.txt", renderdata)
	if err != nil {
		t.Fatal(err)
	}
	err = task.Exec()
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := ioutil.ReadFile(path.Join(tmpdir, "/demo.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "123" {
		t.Fatal(string(bytes))
	}
	bytes, err = ioutil.ReadFile(path.Join(tmpdir, "/demo1.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "data" {
		t.Fatal(string(bytes))
	}
}
