package website

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Jetengine struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Jetengine) ID() string {
	return "github.com/herb-go/herb-go/modules/website/jetengine"
}

func (m *Jetengine) Cmd() string {
	return "jetengine"
}

func (m *Jetengine) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s routers.
Create jetengine module files.
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Jetengine) Desc(a *app.Application) string {
	return "Create jetengine module files."
}
func (m *Jetengine) Group(a *app.Application) string {
	return "Website"
}
func (m *Jetengine) Init(a *app.Application, args *[]string) error {
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
func (m *Jetengine) Question(a *app.Application) error {
	return nil
}
func (m *Jetengine) Exec(a *app.Application, args []string) error {
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

	result, err := tools.FileExists(mp, "app", "website.go")
	if err != nil {
		return err
	}
	if !result {
		err = WebsiteModule.Exec(a, []string{"-s"})
		if err != nil {
			return err
		}

	}

	task := tools.NewTask(filepath.Join(app, "/modules/website/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Jetengine created.\n")
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

func (m *Jetengine) Render(a *app.Application, appPath string, mp string, task *tools.Task) error {
	filesToRender := map[string]string{
		filepath.Join(mp, "views", "jetengine.go"):                        "views/jetengine.go.tmpl",
		filepath.Join("resources", "template.jet", "views.toml"):          "template.jet/views.toml",
		filepath.Join("resources", "template.jet", "layouts", "main.jet"): "template.jet/layouts/main.jet",
		filepath.Join("resources", "template.jet", "views", "index.jet"):  "template.jet/views/index.jet",
	}
	return task.CopyFiles(filesToRender)
}

var JetengineModule = &Jetengine{}

func init() {
	app.Register(JetengineModule)
}
