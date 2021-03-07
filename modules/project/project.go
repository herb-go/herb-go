package project

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

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
	return fmt.Sprintf(help, a.Config.Cmd, a.Config.Cmd)
}

func (m *Project) Desc(a *app.Application) string {
	return "Create new app in given path"
}
func (m *Project) Group(a *app.Application) string {
	return "Application"
}

func (m *Project) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.GoMod, "gomod", false, "use go mod folder struct")
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Project) Exec(a *app.Application, args []string) error {
	mp := "src/vendor/modules"
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		a.PrintModuleHelp(m)
		return nil
	}
	n, err := name.New(true, args...)
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
	var confirmed bool

	err = tools.NewTrueOrFalseQuestion("Do you want to install herb app in ("+appPath+")?").ExecIf(a, true, &confirmed)
	if err != nil {
		return err
	}
	if !confirmed {
		return nil
	}
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/project/resources"), appPath)
	if m.GoMod {
		mp = "src/modules"
		err = task.Render("/skeleton/src/go.mod.example", "/src/go.mod", n)
		if err != nil {
			return err
		}
		err = task.Render("/skeleton/src/vendor/modules/go.mod.example", "/src/modules/go.mod", n)
		if err != nil {
			return err
		}
	} else {
		err = task.Render("/skeleton/src/vendor/modules/go.mod.example", "/src/vendor/modules/go.mod.example", n)
		if err != nil {
			return err
		}
	}
	err = m.createApp(a, appPath, mp, task)
	if err != nil {
		return err
	}

	task.AddJob(func() error {
		a.Printf("App installed in \"%s\"\n", appPath)
		return nil
	})
	return task.Exec()

}

func (m *Project) createApp(a *app.Application, appPath string, mp string, task *tools.Task) error {
	err := task.CopyFiles(AppSkeleton(mp))
	if err != nil {
		return err
	}
	return nil
}

var Module = &Project{}

func init() {
	app.Register(Module)
}
