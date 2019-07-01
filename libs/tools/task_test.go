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
	err = task.Render("/demo1.tmpl", "/output/demo1.txt", renderdata)
	if err != nil {
		t.Fatal(err)
	}
	files := task.ListFiles()
	if len(files) != 2 || files[0] != "/demo.txt" || files[1] != "/output/demo1.txt" {
		t.Fatal(files)
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
	bytes, err = ioutil.ReadFile(path.Join(tmpdir, "/output/demo1.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "data" {
		t.Fatal(string(bytes))
	}
}

func TestTaskFiles(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "herb-go-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(tmpdir)
	}()
	renderdata := map[string]interface{}{
		"data":  "data",
		"data2": "data2",
	}
	task := NewTask(path.Join("./", "testdata"), tmpdir)
	err = task.CopyFiles(map[string]string{"/demo.txt": "/demo.txt", "/demo2.txt": "/demo2.txt"})
	if err != nil {
		t.Fatal(err)
	}
	err = task.RenderFiles(map[string]string{"/demo1.tmpl": "/output/demo1.txt", "/demo3.tmpl": "/output/demo3.txt"}, renderdata)
	if err != nil {
		t.Fatal(err)
	}
	files := task.ListFiles()
	if len(files) != 4 || files[0] != "/demo.txt" || files[1] != "/demo2.txt" || files[2] != "/output/demo1.txt" || files[3] != "/output/demo3.txt" {
		t.Fatal(files)
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
	bytes, err = ioutil.ReadFile(path.Join(tmpdir, "/demo2.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "234" {
		t.Fatal(string(bytes))
	}

	bytes, err = ioutil.ReadFile(path.Join(tmpdir, "/output/demo1.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "data" {
		t.Fatal(string(bytes))
	}
	bytes, err = ioutil.ReadFile(path.Join(tmpdir, "/output/demo3.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "data2" {
		t.Fatal(string(bytes))
	}
}
