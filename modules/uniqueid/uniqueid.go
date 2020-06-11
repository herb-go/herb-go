package uniqueid

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type UniqueID struct {
	app.BasicModule
	SlienceMode bool
}

func (m *UniqueID) ID() string {
	return "github.com/herb-go/herb-go/modules/uniqueid"
}

func (m *UniqueID) Cmd() string {
	return "uniqueid"
}

func (m *UniqueID) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s uniqueid.
Create uniqueid module and config files.
File below will be created:
	config/uniqueid.toml
	system/confg.examples/uniqueid.toml
	src/vendor/modules/app/uniqueid.go
	src/vendor/modules/uniqueid/init.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *UniqueID) Desc(a *app.Application) string {
	return "Create unique id module and config files."
}
func (m *UniqueID) Group(a *app.Application) string {
	return "Misc"
}
func (m *UniqueID) Init(a *app.Application, args *[]string) error {
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
func (m *UniqueID) Question(a *app.Application) error {
	return nil
}
func (m *UniqueID) Exec(a *app.Application, args []string) error {
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

	task := tools.NewTask(filepath.Join(app, "/modules/uniqueid/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("UniqueID   created.\n")
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

func (m *UniqueID) Render(a *app.Application, appPath string, mp string, task *tools.Task) error {

	filesToRender := map[string]string{
		filepath.Join("config", "uniqueid.toml"):                    "uniqueid.toml.tmpl",
		filepath.Join("system", "configskeleton", "uniqueid.toml"): "uniqueid.toml.tmpl",
		filepath.Join(mp, "uniqueid/uniqueid.go"):                   "uniqueid.modules.go.tmpl",
		filepath.Join(mp, "app", "uniqueid.go"):                     "uniqueid.go.tmpl",
	}
	return task.RenderFiles(filesToRender, nil)
}

var UniqueIDModule = &UniqueID{}

func init() {
	app.Register(UniqueIDModule)
}
