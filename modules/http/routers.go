package http

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Routers struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Routers) ID() string {
	return "github.com/herb-go/herb-go/modules/routers"
}

func (m *Routers) Cmd() string {
	return "routers"
}

func (m *Routers) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s routers.
Create routers module files.
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Routers) Desc(a *app.Application) string {
	return "Create routers module files."
}
func (m *Routers) Group(a *app.Application) string {
	return "Web"
}
func (m *Routers) Init(a *app.Application, args *[]string) error {
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
func (m *Routers) Question(a *app.Application) error {
	return nil
}
func (m *Routers) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}

	if len(args) != 0 {
		a.PrintModuleHelp(m)
		return nil
	}
	mp, err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}

	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/http/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Routers created.\n")
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

func (m *Routers) Render(a *app.Application, appPath string, mp string, task *tools.Task) error {
	filesToRender := map[string]string{
		filepath.Join(mp, "routers/routers.go"):             "routers/routers.go.tmpl",
		filepath.Join(mp, "routers/assests.go"):             "routers/assests.go.tmpl",
		filepath.Join(mp, "routers/api.go"):                 "routers/api.go.tmpl",
		filepath.Join(mp, "overseers/router.go"):            "overseers/router.go.tmpl",
		filepath.Join(mp, "overseers/action.go"):            "overseers/action.go.tmpl",
		filepath.Join(mp, "app", "assets.go"):               "assets.go.tmpl",
		filepath.Join("system", "constants", "assets.toml"): "assets.toml.tmpl",
	}
	return task.RenderFiles(filesToRender, nil)
}

var RoutersModule = &Routers{}
