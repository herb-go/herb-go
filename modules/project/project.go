package project

import (
	"github.com/herb-go/util/cli/name"
	"fmt"
	"path"
	"path/filepath"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

const ProjectTypeApp = "app"
const ProjectTypeAPI = "api"
const ProjectTypeWebsite = "website"

const TemplateEngineGoTemple = "tmpl"
const TemplateEngineJet = "jet"

type Project struct {
	app.BasicModule
	ProjectType    string
	TemplateEngine string
	GoMod bool
}

func (m *Project) ID() string {
	return "github.com/herb-go/herb-go/modules/project"
}

func (m *Project) Cmd() string {
	return "new"
}

func (m *Project) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s new <options> [projectname] .
Create new app with given projectname in current folder.

Folder name will filter all host and dir in porject name.
For example,command "%s new github.com/herb-go/newapp" will create folder ./newapp.

`
	return fmt.Sprintf(help, a.Config.Cmd,a.Config.Cmd)
}

func (m *Project) Desc(a *app.Application) string {
	return "Create new app in given path"
}

var projectTypeQuestion = tools.NewQuestion().
	SetDescription("Project type of app").
	AddAnswer("0", "app", ProjectTypeApp).
	AddAnswer("1", "api", ProjectTypeAPI).
	AddAnswer("2", "website", ProjectTypeWebsite).
	SetDefaultKey("1")
var TemplateEngineQuestion = tools.NewQuestion().
	SetDescription("Project engine of website").
	AddAnswer("0", "GO template", TemplateEngineGoTemple).
	AddAnswer("1", "Jet template", TemplateEngineJet).
	SetDefaultKey("0")

func (m *Project) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().StringVar(&m.ProjectType, "type", "", "project type.\"app\",\"api\" or \"website\"")
	m.FlagSet().StringVar(&m.TemplateEngine, "template", "", "website template.\"tmpl\" or \"jet\"")
	m.FlagSet().BoolVar(&m.GoMod, "gomod",false, "use go mod folder struct")
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Project) Question(a *app.Application) error {
	err := projectTypeQuestion.ExecIf(a, m.ProjectType == "", &m.ProjectType)
	if err != nil {
		return err
	}
	err = TemplateEngineQuestion.ExecIf(a, m.TemplateEngine == "" && m.ProjectType == ProjectTypeWebsite, &m.TemplateEngine)
	if err != nil {
		return err
	}
	return nil
}
func (m *Project) Exec(a *app.Application, args []string) error {
	mp:="src/vendor/modules"
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		a.PrintModuleHelp(m)
		return nil
	}
	n,err:=name.New(true,args...)
		if err != nil {
		return err
	}
	appPath := path.Join(a.Cwd, n.Lower)
	result, err := tools.FileExists(appPath)
	if err != nil {
		return err
	} 
	if result {
		return fmt.Errorf("\"%s\" exists.Create app fail", appPath)
	}
	err = tools.ErrorIfStringFieldNotInList("type", m.ProjectType, "", ProjectTypeAPI, ProjectTypeApp, ProjectTypeWebsite)
	if err != nil {
		return err
	}
	err = tools.ErrorIfStringFieldNotInList("template", m.TemplateEngine, "", TemplateEngineGoTemple, TemplateEngineGoTemple)
	if err != nil {
		return err
	}
	err = m.Question(a)
	if err != nil {
		return err
	}

	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/project/resources"), appPath)
	if m.GoMod {
		mp="src/modules"
		err=task.Render("/skeleton/src/go.mod.example","/src/go.mod",n)
		if err != nil {
			return err
		}
	}
	err = m.createApp(a, appPath, mp,task)
	if err != nil {
		return err
	}

	if m.ProjectType == ProjectTypeAPI || m.ProjectType == ProjectTypeWebsite {
		err = m.createHTTP(a, appPath,mp, task)
		if err != nil {
			return err
		}
	}
	if m.ProjectType == ProjectTypeWebsite {
		err = m.createWebsite(a, appPath,mp, task)
		if err != nil {
			return err
		}
		if m.TemplateEngine == TemplateEngineJet {
			err = m.createJetEngine(a, appPath,mp, task)
			if err != nil {
				return err
			}
		}
		if m.TemplateEngine == TemplateEngineGoTemple {
			err = m.createTmplEngine(a, appPath, mp,task)
			if err != nil {
				return err
			}
		}
	}

	task.AddJob(func() error {
		a.Printf("App installed in \"%s\"\n", appPath)
		return nil
	})
	return task.Exec()

}

func (m *Project)  createApp(a *app.Application, appPath string, mp string,task *tools.Task) error {
	err := task.CopyFiles(AppSkeleton(mp))
	if err != nil {
		return err
	}
	return nil
}

func (m *Project)  createHTTP(a *app.Application, appPath string,mp string, task *tools.Task) error {
	err := task.CopyFiles(HTTPSkeleton(mp))
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
func (m *Project)  createWebsite(a *app.Application, appPath string, mp string,task *tools.Task) error {
	err := task.CopyFiles(WebsiteSkeleton(mp))
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

func (m *Project)  createJetEngine(a *app.Application, appPath string, mp string,task *tools.Task) error {
	err := task.CopyFiles(JetEngineSkeleton(mp))
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

func (m *Project)  createTmplEngine(a *app.Application, appPath string, mp string,task *tools.Task) error {
	err := task.CopyFiles(TmplEngineSkeleton(mp))
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
