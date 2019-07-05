package config

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Constants struct {
	app.BasicModule
	Watch       bool
	AutoConfirm bool
}

func (m *Constants) ID() string {
	return "github.com/herb-go/herb-go/modules/config.constants"
}

func (m *Constants) Cmd() string {
	return "constants"
}

func (m *Constants) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s constants [name].
Create new constant file and code.
File below will be created:
	system/constant/[name].toml
	src/vendor/modules/app/[name].go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Constants) Desc(a *app.Application) string {
	return "Create new constants file and code"
}
func (m *Constants) Group(a *app.Application) string {
	return "Config"
}
func (m *Constants) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.AutoConfirm, "y", false, "Whether auto confirm")

	m.FlagSet().BoolVar(&m.Watch, "watch", false, "Whether auto reload constants  after file changed")
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Constants) Question(a *app.Application) error {
	return nil
}
func (m *Constants) Exec(a *app.Application, args []string) error {
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
	mp,err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}

	task := tools.NewTask(filepath.Join(app, "/modules/config/resources"), a.Cwd)

	err = m.Render(a, a.Cwd,mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Config  \"%s\" created.\n", n.LowerWithParentDotSeparated)
		return nil
	})
	err = task.ErrosIfAnyFileExists()
	if err != nil {
		return err
	}
	ok, err := task.ConfirmIf(a, !m.AutoConfirm)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return task.Exec()

}

func (m *Constants) Render(a *app.Application, appPath string,mp string, task *tools.Task, n *name.Name) error {
	var configgopath string
	if m.Watch {
		configgopath = "constantswatch.go.tmpl"
	} else {
		configgopath = "constants.go.tmpl"
	}
	filesToRender := map[string]string{
		filepath.Join("system", "constants", n.LowerWithParentDotSeparated+".toml"):           "config.toml.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"): configgopath,
	}
	return task.RenderFiles(filesToRender, n)
}

var ConstantsModule = &Constants{}

func init() {
	app.Register(ConstantsModule)
}
