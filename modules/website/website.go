package website

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/http"
	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Website struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Website) ID() string {
	return "github.com/herb-go/herb-go/modules/website"
}

func (m *Website) Cmd() string {
	return "website"
}

func (m *Website) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s routers.
Create website module files.
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Website) Desc(a *app.Application) string {
	return "Create website module files."
}
func (m *Website) Group(a *app.Application) string {
	return "Website"
}
func (m *Website) Init(a *app.Application, args *[]string) error {
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
func (m *Website) Question(a *app.Application) error {
	return nil
}
func (m *Website) Exec(a *app.Application, args []string) error {
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

	result, err := tools.FileExists(mp, "middlewares", "middlewares.go")
	if err != nil {
		return err
	}
	if !result {
		err = http.MiddlewaresModule.Exec(a, []string{"-s"})
		if err != nil {
			return err
		}

	}
	result, err = tools.FileExists(mp, "routers", "routers.go")
	if err != nil {
		return err
	}
	if !result {
		err = http.RoutersModule.Exec(a, []string{"-s"})
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
		a.Printf("Website created.\n")
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

func (m *Website) Render(a *app.Application, appPath string, mp string, task *tools.Task) error {
	filesToRender := map[string]string{
		filepath.Join(mp, "app", "website.go"):                               "website.go.tmpl",
		filepath.Join(mp, "actions", "website.go"):                           "actions/website.go.tmpl",
		filepath.Join(mp, "middlewares", "errorpages.go"):                    "middlewares/errorpages.go.tmpl",
		filepath.Join(mp, "views", "init.go"):                                "views/init.go.tmpl",
		filepath.Join(mp, "views", "views.go"):                               "views/views.go.tmpl",
		filepath.Join("system", "defaultconfig", "website.toml"):             "website.toml.tmpl",
		filepath.Join("system", "defaultconfig", "examples", "website.toml"): "website.toml.tmpl",
		filepath.Join("system", "defaultconfig", "examples", "website.toml"): "website.toml.tmpl",

		filepath.Join("resources", "errorpages", "404.html"): "errorpages/404.html",
		filepath.Join("resources", "errorpages", "500.html"): "errorpages/500.html",
	}
	return task.RenderFiles(filesToRender, nil)
}

var WebsiteModule = &Website{}

func init() {
	app.Register(WebsiteModule)
}
