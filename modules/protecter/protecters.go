package protecter

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/httpinfo"
	"github.com/herb-go/herb-go/modules/overseers"
	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Protecters struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Protecters) ID() string {
	return "github.com/herb-go/herb-go/modules/protecter"
}

func (m *Protecters) Cmd() string {
	return "protecters"
}

func (m *Protecters) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s protecters <name>.
Create protecters module and config files.
Default name is "protecter".
File below will be created:
	config/<name>.toml
	system/configskeleton/<name>.toml
	src/vendor/modules/app/<name>.go
	src/vendor/modules/<name>/init.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Protecters) Desc(a *app.Application) string {
	return "Create protecters module and config files."
}
func (m *Protecters) Group(a *app.Application) string {
	return "Auth"
}
func (m *Protecters) Init(a *app.Application, args *[]string) error {
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
func (m *Protecters) Question(a *app.Application) error {
	return nil
}
func (m *Protecters) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}

	if len(args) != 0 {
		a.PrintModuleHelp(m)
		return nil
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
	err = overseers.OverseerModule.Exec(a, []string{"-s", "authenticator"})
	if err != nil {
		return err
	}
	result, err := tools.FileExists(mp, "httpinfo", "init.go")
	if err != nil {
		return err
	}
	if !result {
		err = httpinfo.HTTPInfoModule.Exec(a, []string{"-s"})
		if err != nil {
			return err
		}

	}
	task := tools.NewTask(filepath.Join(app, "/modules/herbmodules/protecter/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Protecters  created.\n")
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

func (m *Protecters) Render(a *app.Application, appPath string, mp string, task *tools.Task) error {
	err := tools.CopyIfNotExist(filepath.Join(task.SrcFolder, "hired.authenticator.go.tmpl"), mp, "hired", "authenticator.go")
	if err != nil {
		return err
	}
	filesToRender := map[string]string{
		filepath.Join("config", "protecters.toml"):                   "protecters.toml.tmpl",
		filepath.Join("system", "configskeleton", "protecters.toml"): "protecters.toml.tmpl",
		filepath.Join(mp, "protecters", "init.go"):                   "protecters.go.tmpl",
		filepath.Join(mp, "app", "protecters.go"):                    "app.protecters.go.tmpl",
		filepath.Join(mp, "protecters.go"):                           "modules.go.tmpl",
	}
	return task.RenderFiles(filesToRender, nil)
}

var ProtectersModule = &Protecters{}

func init() {
	app.Register(ProtectersModule)
}
