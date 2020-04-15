package overseers

import (
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

var OverSeerInitFuncs = map[string]func(a *app.Application, appPath string, mp string, slienceMode bool) error{
	"action": project.InitActionOverseer,
	"cache":  newInitFunc("cache.go", "cache.go", "Cache"),
}

func newInitFunc(dst string, src string, name string) func(a *app.Application, appPath string, mp string, slienceMode bool) error {
	return func(a *app.Application, appPath string, mp string, slienceMode bool) error {
		var result bool
		var err error
		app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
		if err != nil {
			return err
		}
		task := tools.NewTask(filepath.Join(app, "/modules/overseers/resources"), appPath)
		filesToRender := map[string]string{}
		result, err = tools.FileExists(filepath.Join(appPath, mp, "overseers", dst))
		if err != nil {
			return err
		}
		if !result {
			filesToRender[filepath.Join(mp, "overseers", dst)] = src
		}
		if len(filesToRender) == 0 {
			return nil
		}
		task.AddJob(func() error {
			a.Printf("%s overseer inited.\n", name)
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
}
