package tools

import (
	"bytes"
	"io/ioutil"
	"path"
	"text/template"
)

func NewTask(srcfolder string, targetfolder string) *Task {
	return &Task{
		SrcFolder:    srcfolder,
		TargetFolder: targetfolder,
		Files:        map[string][]byte{},
	}
}

type Task struct {
	SrcFolder    string
	TargetFolder string
	Files        map[string][]byte
}

func (t *Task) Copy(src string, target string) error {
	bs, err := ioutil.ReadFile(path.Join(t.SrcFolder, src))
	if err != nil {
		return err
	}
	t.Files[path.Join(t.TargetFolder, target)] = bs
	return nil
}

func (t *Task) Render(src string, target string, data interface{}) error {
	tmpl, err := template.ParseFiles(path.Join(t.SrcFolder, src))
	if err != nil {
		return err
	}
	bs := []byte{}
	err = tmpl.Execute(bytes.NewBuffer(bs), data)
	if err != nil {
		return err
	}
	t.Files[path.Join(t.TargetFolder, target)] = bs
	return nil
}
