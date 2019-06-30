package tools

import (
	"bytes"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/herb-go/util"
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
	t.Files[target] = bs
	return nil
}

func (t *Task) Render(src string, target string, data interface{}) error {
	tmpl, err := template.ParseFiles(path.Join(t.SrcFolder, src))
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, data)
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadAll(buf)
	if err != nil {
		return err
	}
	t.Files[target] = bs
	return nil
}
func (t *Task) ListFiles() []string {
	result := []string{}
	for k := range t.Files {
		result = append(result, k)
	}
	return result
}
func (t *Task) Exec() error {
	for k := range t.Files {
		target := path.Join(t.TargetFolder, k)
		err := ioutil.WriteFile(target, t.Files[k], util.DefaultFileMode)
		if err != nil {
			return err
		}
	}
	return nil
}
