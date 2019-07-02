package project

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/herb-go/herb-go/app"
	"github.com/herb-go/herb-go/libs/tools"
)

type Project struct {
	app.BasicModule
}

func (m *Project) ID() string {
	return "github.com/herb-go/herb-go/modules/project"
}

func (m *Project) Cmd() string {
	return "new"
}

func (m *Project) Help() string {
	return ""
}

func (m *Project) Desc() string {
	return ""
}
func (m *Project) Exec(a *app.Application, args []string) error {
	appPath := path.Join(a.Cwd, args[0])
	result, err := tools.FileExists(appPath)
	if err != nil {
		return err
	}
	if result {
		return fmt.Errorf("\"%s\" exists.Create app fail", appPath)
	}
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/project/resources"), appPath)
	err = createApp(a, appPath, task)
	if err != nil {
		return err
	}
	err = createHTTP(a, appPath, task)
	if err != nil {
		return err
	}
	err = createWebsite(a, appPath, task)
	if err != nil {
		return err
	}
	// err = createJetEngine(a, appPath, task)
	// if err != nil {
	// 	return err
	// }
	err = createTmplEngine(a, appPath, task)
	if err != nil {
		return err
	}
	err = task.Exec()
	if err != nil {
		return err
	}
	a.Printf("App installed in \"%s\"\n", appPath)
	return nil

}

func createApp(a *app.Application, appPath string, task *tools.Task) error {
	err := task.CopyFiles(AppSkeleton)
	if err != nil {
		return err
	}
	return nil
}

func createHTTP(a *app.Application, appPath string, task *tools.Task) error {
	err := task.CopyFiles(HTTPSkeleton)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		_, err := tools.ReplaceLine(filepath.Join(appPath, "/src/main.go"),
			`//Replace next line "panic(errFuncWhenRunFuncNotRewrited)" with your own app run function`,
			"	//Run app as http server.",
		)
		if err != nil {
			return err
		}
		_, err = tools.ReplaceLine(filepath.Join(appPath, "/src/main.go"),
			`panic(errFuncWhenRunFuncNotRewrited)`,
			"	RunHTTP()",
		)
		return err
	})
	return nil
}
func createWebsite(a *app.Application, appPath string, task *tools.Task) error {
	err := task.CopyFiles(WebsiteSkeleton)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		_, err := tools.ReplaceLine(filepath.Join(appPath, "/src/vendor/modules/routers/routers.go"),
			`//"modules/actions"`,
			`	"modules/actions"`,
		)
		if err != nil {
			return err
		}
		_, err = tools.ReplaceLine(filepath.Join(appPath, "/src/vendor/modules/routers/routers.go"),
			`//var RouterHTML = newHTMLRouter()`,
			"	var RouterHTML = newHTMLRouter()",
		)
		if err != nil {
			return err
		}
		_, err = tools.ReplaceLine(filepath.Join(appPath, "/src/vendor/modules/routers/routers.go"),
			`//Router.StripPrefix("/page").Use(HTMLMiddlewares()...).Handle(RouterHTML)`,
			"	Router.StripPrefix(\"/page\").\r\n		Use(HTMLMiddlewares()...).\r\n		Handle(RouterHTML)",
		)

		return err
	})
	return nil
}

func createJetEngine(a *app.Application, appPath string, task *tools.Task) error {
	err := task.CopyFiles(JetEngineSkeleton)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		_, err := tools.ReplaceLine(filepath.Join(appPath, "/src/vendor/modules/views/init.go"),
			"var ViewsInitiator func()",
			"var ViewsInitiator = initJetViews",
		)
		if err != nil {
			return err
		}
		_, err = tools.ReplaceLine(filepath.Join(appPath, "/src/vendor/modules/routers/routers.go"),
			`//Router.GET("/").Use(HTMLMiddlewares()...).HandleFunc(actions.IndexAction)`,
			"	Router.GET(\"/\").\n		Use(HTMLMiddlewares()...).\n		HandleFunc(actions.IndexAction)",
		)

		return err
	})
	return nil
}

func createTmplEngine(a *app.Application, appPath string, task *tools.Task) error {
	err := task.CopyFiles(TmplEngineSkeleton)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		_, err := tools.ReplaceLine(filepath.Join(appPath, "/src/vendor/modules/views/init.go"),
			"var ViewsInitiator func()",
			"var ViewsInitiator = initTmplViews",
		)
		if err != nil {
			return err
		}
		_, err = tools.ReplaceLine(filepath.Join(appPath, "/src/vendor/modules/routers/routers.go"),
			`//Router.GET("/").Use(HTMLMiddlewares()...).HandleFunc(actions.IndexAction)`,
			"	Router.GET(\"/\").\n		Use(HTMLMiddlewares()...).\n		HandleFunc(actions.IndexAction)",
		)

		return err
	})
	return nil
}

var Module = &Project{}

func init() {
	app.Register(Module)
}
