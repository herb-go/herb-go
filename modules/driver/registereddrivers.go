package driver

import (
	"path/filepath"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

func GetRegisteredDrivers() []string {
	var result = []string{}
	for k := range DriverInitFuncs {
		result = append(result, k)
	}
	return result
}

var DriverInitFuncs = map[string]func(a *app.Application, appPath string, mp string, slienceMode bool) error{
	"kvdb":         newInitFunc([]string{"kvdb.go", "kvdb.go"}, "Kvdb"),
	"texttemplate": newInitFunc([]string{"texttemplate.go", "texttemplate.go"}, "TextTemplate"),
}

//files:[]{"dst","src","dst2","src2"}
func newInitFunc(files []string, name string) func(a *app.Application, appPath string, mp string, slienceMode bool) error {
	return func(a *app.Application, appPath string, mp string, slienceMode bool) error {
		var result bool
		var err error
		app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
		if err != nil {
			return err
		}
		task := tools.NewTask(filepath.Join(app, "/modules/driver/resources"), appPath)
		filesToRender := map[string]string{}
		for i := 0; i < len(files); i = i + 2 {
			result, err = tools.FileExists(filepath.Join(appPath, mp, "drivers", files[i]))
			if err != nil {
				return err
			}
			if !result {
				filesToRender[filepath.Join(mp, "drivers", files[i])] = files[i+1]
			}
			if len(filesToRender) == 0 {
				return nil
			}
		}
		task.AddJob(func() error {
			a.Printf("%s driver added.\n", name)
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
