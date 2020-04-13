package worker

import (
	"path/filepath"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

func InitWorkers(a *app.Application, appPath string, mp string, slienceMode bool) error {
	var result bool
	var err error
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/worker/resources"), a.Cwd)
	filesToRender := map[string]string{}
	result, err = tools.FileExists(filepath.Join(mp, "workers", "init.go"))
	if err != nil {
		return err
	}
	if !result {
		filesToRender[filepath.Join(mp, "workers", "init.go")] = "workers_init.go.tmpl"
	}
	result, err = tools.FileExists(filepath.Join(mp, "workers", "overseers", "init.go"))
	if err != nil {
		return err
	}
	if !result {
		filesToRender[filepath.Join(mp, "workers", "overseers", "init.go")] = "overseers_init.go.tmpl"
	}
	if len(filesToRender) == 0 {
		return nil
	}
	task.AddJob(func() error {
		a.Printf("Workers inited.\n")
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
