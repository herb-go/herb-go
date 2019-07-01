package tools

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/herb-go/herb-go/app"
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
func (t *Task) CopyFiles(files map[string]string) error {
	for k := range files {
		err := t.Copy(k, files[k])
		if err != nil {
			return err
		}
	}
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

func (t *Task) RenderFiles(files map[string]string, data interface{}) error {
	for k := range files {
		err := t.Render(k, files[k], data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) ListFiles() []string {
	result := []string{}
	for k := range t.Files {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}
func (t *Task) Exec() error {
	for k := range t.Files {
		target := filepath.Join(t.TargetFolder, k)
		targetdir := filepath.Dir(target)
		if targetdir != "" {
			os.MkdirAll(targetdir, util.DefaultFolderMode)
		}
		err := ioutil.WriteFile(target, t.Files[k], util.DefaultFileMode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) ConfirmIf(a *app.Application, conditon bool) (bool, error) {
	if !conditon {
		return true, nil
	}
	var result bool
	files := t.ListFiles()
	if len(files) == 0 {
		return true, nil
	}
	a.Println("Below files will be installed:")
	for _, v := range files {
		a.Println(v)
	}
	q := NewQuestion().
		SetDescription("Install all files?").
		AddAnswer("y", "Yes", true).
		AddAnswer("n", "No", false)
	err := q.ExecIf(a, true, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}
