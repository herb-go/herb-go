package config

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
	"github.com/herb-go/util/cli/name"
)

type DefaultConfig struct {
	app.BasicModule
	Watch       bool
	SlienceMode bool
}

func (m *DefaultConfig) ID() string {
	return "github.com/herb-go/herb-go/modules/config"
}

func (m *DefaultConfig) Cmd() string {
	return "defaultconfig"
}

func (m *DefaultConfig) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s defaultconfig [name].
Create new defaultconfig file and code.
File below will be created:
	system/defaultconfig/[name].toml
	system/defaultconfig/examples/[name].toml
	src/vendor/modules/app/[name].go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *DefaultConfig) Desc(a *app.Application) string {
	return "Create new default config file and code"
}
func (m *DefaultConfig) Group(a *app.Application) string {
	return "Config"
}
func (m *DefaultConfig) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.SlienceMode, "s", false, "Slience mode")

	m.FlagSet().BoolVar(&m.Watch, "watch", false, "Whether auto reload config after file changed")

	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *DefaultConfig) Question(a *app.Application) error {
	return nil
}
func (m *DefaultConfig) Exec(a *app.Application, args []string) error {
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

	task := tools.NewTask(filepath.Join(app, "/modules/config/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
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
	ok, err := task.ConfirmIf(a, !m.SlienceMode)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return task.Exec()

}

func (m *DefaultConfig) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {
	var configgopath string
	if m.Watch {
		configgopath = "configwatch.go.tmpl"
	} else {
		configgopath = "config.go.tmpl"
	}
	filesToRender := map[string]string{
		filepath.Join("system", "defaultconfig", n.LowerWithParentDotSeparated+".toml"):             "config.toml.tmpl",
		filepath.Join("system", "defaultconfig", "examples", n.LowerWithParentDotSeparated+".toml"): "config.toml.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"):                               configgopath,
	}
	return task.RenderFiles(filesToRender, n)
}

var DefaultConfigModule = &DefaultConfig{}

func init() {
	app.Register(DefaultConfigModule)
}
