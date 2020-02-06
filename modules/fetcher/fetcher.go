package fetcher

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type API struct {
	app.BasicModule
	SlienceMode bool
}

func (m *API) ID() string {
	return "github.com/herb-go/herb-go/modules/fetcher"
}

func (m *API) Cmd() string {
	return "fetcher"
}

func (m *API) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s fetcher [name].
Create new fetcher config and go file.
File below will be created:
	config/[name].toml
	system/config.examples/[name].toml
	src/vendor/modules/app/[name].go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *API) Desc(a *app.Application) string {
	return "Create config file and code to do http request"
}
func (m *API) Group(a *app.Application) string {
	return "Third part"
}
func (m *API) Init(a *app.Application, args *[]string) error {
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
func (m *API) Question(a *app.Application) error {
	return nil
}
func (m *API) Exec(a *app.Application, args []string) error {
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

	task := tools.NewTask(filepath.Join(app, "/modules/fetcher/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Fetcher  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *API) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {

	filesToRender := map[string]string{
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):                    "fetcher.toml.tmpl",
		filepath.Join("system", "config.examples", n.LowerWithParentDotSeparated+".toml"): "fetcher.toml.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"):                     "fetcher.go.tmpl",
	}
	return task.RenderFiles(filesToRender, n)
}

var APIModule = &API{}

func init() {
	app.Register(APIModule)
}
