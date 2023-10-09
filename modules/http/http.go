package http

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type HTTP struct {
	app.BasicModule
	SlienceMode bool
}

func (m *HTTP) ID() string {
	return "github.com/herb-go/herb-go/modules/http"
}

func (m *HTTP) Cmd() string {
	return "http"
}

func (m *HTTP) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s http <name>.
Create http module and config files.
Default name is "http".
File below will be created:
	config/<name>.toml
	system/configskeleton/<name>.toml
	src/modules/app/<name>.go
	src/modules/<name>/init.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *HTTP) Desc(a *app.Application) string {
	return "Create new http code and config"
}

func (m *HTTP) Group(a *app.Application) string {
	return "Web"
}

func (m *HTTP) Init(a *app.Application, args *[]string) error {
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
func (m *HTTP) Question(a *app.Application) error {
	return nil
}
func (m *HTTP) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name
	if len(args) == 0 {
		a.Println("No http module name given.\"http\" is used")
		n, err = name.New(true, "http")
	} else {
		n, err = name.New(true, args...)
	}
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
	result, err := tools.FileExists(mp, "middlewares", "middlewares.go")
	if err != nil {
		return err
	}
	if !result {
		err = MiddlewaresModule.Exec(a, []string{"-s"})
		if err != nil {
			return err
		}

	}
	result, err = tools.FileExists(mp, "routers", "routers.go")
	if err != nil {
		return err
	}
	if !result {
		err = RoutersModule.Exec(a, []string{"-s"})
		if err != nil {
			return err
		}

	}
	task := tools.NewTask(filepath.Join(app, "/modules/http/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("HTTP  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *HTTP) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {

	filesToRender := map[string]string{
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):                   "http.toml.tmpl",
		filepath.Join("system", "configskeleton", n.LowerWithParentDotSeparated+".toml"): "http.toml.tmpl",
		filepath.Join(mp, n.LowerWithParentDotSeparated+".go"):                           "http.modules.tmpl",
		filepath.Join(mp, n.LowerPath("init.go")):                                        "http.go.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"):                    "app.http.go.tmpl",
	}
	return task.RenderFiles(filesToRender, n)
}

var Module = &HTTP{}

func init() {
	app.Register(Module)
}
