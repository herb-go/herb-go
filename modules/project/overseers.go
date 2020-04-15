package project

import (
	"path/filepath"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

func InitOverseers(a *app.Application, appPath string, mp string, slienceMode bool) error {
	var result bool
	var err error
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/project/resources/skeleton/src/vendor/modules/overseers"), appPath)
	filesToRender := map[string]string{}
	result, err = tools.FileExists(filepath.Join(appPath, mp, "overseers", "init.go"))
	if err != nil {
		return err
	}
	if !result {
		filesToRender[filepath.Join(mp, "overseers", "init.go")] = "init.go"
	}
	if len(filesToRender) == 0 {
		return nil
	}
	task.AddJob(func() error {
		a.Printf("Overseers inited.\n")
		return nil
	})
	err = task.RenderFiles(filesToRender, nil)
	if err != nil {
		return err
	}
	ok, err := task.ConfirmIf(a, !slienceMode)
	if err != nil {
		return err
	}
	if !ok {
		return tools.ErrUserCanceled
	}
	return task.Exec()
}

func InitActionOverseer(a *app.Application, appPath string, mp string, slienceMode bool) error {
	var result bool
	var err error
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/project/resources/skeleton/src/vendor/modules/overseers"), appPath)
	filesToRender := map[string]string{}
	result, err = tools.FileExists(filepath.Join(appPath, mp, "overseers", "action.go"))
	if err != nil {
		return err
	}
	if !result {
		filesToRender[filepath.Join(mp, "overseers", "action.go")] = "action.go"
	}
	if len(filesToRender) == 0 {
		return nil
	}
	task.AddJob(func() error {
		a.Printf("Action overseer inited.\n")
		return nil
	})
	err = task.CopyFiles(filesToRender)
	if err != nil {
		return err
	}
	ok, err := task.ConfirmIf(a, !slienceMode)
	if err != nil {
		return err
	}
	if !ok {
		return tools.ErrUserCanceled
	}
	return task.Exec()
}
