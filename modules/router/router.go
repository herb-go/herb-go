package router

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/http"
	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Router struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Router) ID() string {
	return "github.com/herb-go/herb-go/modules/router"
}

func (m *Router) Cmd() string {
	return "router"
}

func (m *Router) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s router [name].
Create new router.
File below will be created:
	routers/[name].go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Router) Desc(a *app.Application) string {
	return "Create new router"
}
func (m *Router) Group(a *app.Application) string {
	return "Web"
}
func (m *Router) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.SlienceMode, "s", false, "Slience mode")
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Router) Question(a *app.Application) error {
	return nil
}
func (m *Router) Exec(a *app.Application, args []string) error {
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
	mp, err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}
	result, err := tools.FileExists(mp, "routers", "routers.go")
	if err != nil {
		return err
	}
	if !result {
		err = http.RoutersModule.Exec(a, []string{"-s"})
		if err != nil {
			return err
		}

	}

	task := tools.NewTask(filepath.Join(app, "/modules/router/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Router  \"%s\" created.\n", n.LowerWithParentDotSeparated)
		return nil
	})
	err = task.ErrosIfAnyFileExists()
	if err != nil {
		return err
	}
	ok, err := task.ConfirmIf(a, !m.SlienceMode)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return task.Exec()

}

func (m *Router) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {

	filesToRender := map[string]string{
		filepath.Join(mp, "routers", n.LowerWithParentDotSeparated+".go"): "router.go.tmpl",
	}
	return task.RenderFiles(filesToRender, n)
}

var RouterModule = &Router{}

func init() {
	app.Register(RouterModule)
}
