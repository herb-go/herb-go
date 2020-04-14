package overseers

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
	task := tools.NewTask(filepath.Join(app, "/modules/overseer/resources"), a.Cwd)
	filesToRender := map[string]string{}
	result, err = tools.FileExists(filepath.Join(mp, "overseers", "init.go"))
	if err != nil {
		return err
	}
	if !result {
		filesToRender[filepath.Join(mp, "overseers", "init.go")] = "overseers_init.go.tmpl"
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
