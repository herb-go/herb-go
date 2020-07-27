package website

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Tmplengine struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Tmplengine) ID() string {
	return "github.com/herb-go/herb-go/modules/website/tmplengine"
}

func (m *Tmplengine) Cmd() string {
	return "tmplengine"
}

func (m *Tmplengine) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s routers.
Create tmplengine module files.
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Tmplengine) Desc(a *app.Application) string {
	return "Create tmplengine module files."
}
func (m *Tmplengine) Group(a *app.Application) string {
	return "Website"
}
func (m *Tmplengine) Init(a *app.Application, args *[]string) error {
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
func (m *Tmplengine) Question(a *app.Application) error {
	return nil
}
func (m *Tmplengine) Exec(a *app.Application, args []string) error {
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
		a.Printf("Tmplengine created.\n")
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

func (m *Tmplengine) Render(a *app.Application, appPath string, mp string, task *tools.Task) error {
	filesToRender := map[string]string{
		filepath.Join(mp, "views", "tmplengine.go"):                         "views/tmplengine.go.tmpl",
		filepath.Join("resources", "template.tmpl", "views.toml"):           "template.tmpl/views.toml",
		filepath.Join("resources", "template.tmpl", "layouts", "main.tmpl"): "template.tmpl/layouts/main.tmpl",
		filepath.Join("resources", "template.tmpl", "views", "index.tmpl"):  "template.tmpl/views/index.tmpl",
	}
	return task.CopyFiles(filesToRender)
}

var TmplengineModule = &Tmplengine{}

func init() {
	app.Register(TmplengineModule)
}
