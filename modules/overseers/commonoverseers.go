package overseers

import (
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

var OverSeerInitFuncs = map[string]func(a *app.Application, appPath string, mp string, slienceMode bool) error{
	"action":        project.InitActionOverseer,
	"cache":         newInitFunc([]string{"cache.go", "cache.go", "cacheproxy.go", "cacheproxy.go"}, "Cache"),
	"member":        newInitFunc([]string{"member.go", "member.go", "memberdirectivefactory.go", "memberdirectivefactory.go"}, "Member"),
	"authenticator": newInitFunc([]string{"authenticatorfactory.go", "authenticatorfactory.go"}, "Authenticator"),
	"database":      newInitFunc([]string{"database.go", "database.go"}, "Database"),
	"identifier":    newInitFunc([]string{"identifier.go", "identifier.go"}, "Identifier"),
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
		task := tools.NewTask(filepath.Join(app, "/modules/overseers/resources"), appPath)
		filesToRender := map[string]string{}
		for i := 0; i < len(files); i = i + 2 {
			result, err = tools.FileExists(filepath.Join(appPath, mp, "overseers", files[i]))
			if err != nil {
				return err
			}
			if !result {
				filesToRender[filepath.Join(mp, "overseers", files[i])] = files[i+1]
			}
			if len(filesToRender) == 0 {
				return nil
			}
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
