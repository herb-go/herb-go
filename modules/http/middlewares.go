package http

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Middlewares struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Middlewares) ID() string {
	return "github.com/herb-go/herb-go/modules/middlewares"
}

func (m *Middlewares) Cmd() string {
	return "middlewares"
}

func (m *Middlewares) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s middlewares.
Create middlewares module files.
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Middlewares) Desc(a *app.Application) string {
	return "Create middlewares module files."
}
func (m *Middlewares) Group(a *app.Application) string {
	return "Web"
}
func (m *Middlewares) Init(a *app.Application, args *[]string) error {
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
func (m *Middlewares) Question(a *app.Application) error {
	return nil
}
func (m *Middlewares) Exec(a *app.Application, args []string) error {
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
		a.Printf("Middlewares created.\n")
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

func (m *Middlewares) Render(a *app.Application, appPath string, mp string, task *tools.Task) error {
	filesToRender := map[string]string{
		filepath.Join(mp, "middlewares/middlewares.go"):                   "middlewares/middlewares.go.tmpl",
		filepath.Join(mp, "middlewares/csrf.go"):                          "middlewares/csrf.go.tmpl",
		filepath.Join(mp, "middlewares/factory.go"):                       "middlewares/factory.go.tmpl",
		filepath.Join(mp, "overseers/middleware.go"):                      "overseers/middleware.go.tmpl",
		filepath.Join(mp, "app", "csrf.go"):                               "csrf.go.tmpl",
		filepath.Join("config", "csrf.toml"):                              "csrf.toml.tmpl",
		filepath.Join("system", "defaultconfig", "csrf.toml"):             "csrf.toml.tmpl",
		filepath.Join("system", "defaultconfig", "examples", "csrf.toml"): "csrf.toml.tmpl",
	}
	return task.RenderFiles(filesToRender, nil)
}

var MiddlewaresModule = &Middlewares{}
