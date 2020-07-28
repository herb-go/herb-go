package cors

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type CORS struct {
	app.BasicModule
	SlienceMode bool
}

func (m *CORS) ID() string {
	return "github.com/herb-go/herb-go/modules/cors"
}

func (m *CORS) Cmd() string {
	return "cors"
}

func (m *CORS) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s cors <name>.
Create cors module and config files.
Default name is "cors".
File below will be created:
	config/<name>.toml
	system/configskeleton/<name>.toml
	src/vendor/modules/app/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *CORS) Desc(a *app.Application) string {
	return "Create new cors code and config"
}

func (m *CORS) Group(a *app.Application) string {
	return "Web"
}

func (m *CORS) Init(a *app.Application, args *[]string) error {
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
func (m *CORS) Question(a *app.Application) error {
	return nil
}
func (m *CORS) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name
	if len(args) == 0 {
		a.Println("No cors module name given.\"cors\" is used")
		n, err = name.New(true, "cors")
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

	task := tools.NewTask(filepath.Join(app, "/modules/cors/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("CORS  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *CORS) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {

	filesToRender := map[string]string{
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):                   "cors.toml.tmpl",
		filepath.Join("system", "configskeleton", n.LowerWithParentDotSeparated+".toml"): "cors.toml.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"):                    "cors.go.tmpl",
	}
	return task.RenderFiles(filesToRender, n)
}

var Module = &CORS{}

func init() {
	app.Register(Module)
}
