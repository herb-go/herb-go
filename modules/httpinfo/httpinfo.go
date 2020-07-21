package httpinfo

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type HTTPInfo struct {
	app.BasicModule
	SlienceMode bool
}

func (m *HTTPInfo) ID() string {
	return "github.com/herb-go/herb-go/modules/httpinfo"
}

func (m *HTTPInfo) Cmd() string {
	return "httpinfo"
}

func (m *HTTPInfo) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s httpinfo.
Create httpinfo config and code.
File below will be created:
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *HTTPInfo) Desc(a *app.Application) string {
	return "Create httpinfo config and code"
}
func (m *HTTPInfo) Group(a *app.Application) string {
	return "Web"
}
func (m *HTTPInfo) Init(a *app.Application, args *[]string) error {
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
func (m *HTTPInfo) Question(a *app.Application) error {
	return nil
}
func (m *HTTPInfo) Exec(a *app.Application, args []string) error {
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

	task := tools.NewTask(filepath.Join(app, "/modules/httpinfo/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("HTTPInfo module  created.\n")
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

func (m *HTTPInfo) Render(a *app.Application, appPath string, mp string, task *tools.Task) error {

	filesToRender := map[string]string{
		filepath.Join(mp, "drivers", "httpinfo.go"):                 "driver.go.tmpl",
		filepath.Join(mp, "httpinfo.go"):                            "module.go.tmpl",
		filepath.Join(mp, "httpinfo", "init.go"):                    "httpinfo.init.go.tmpl",
		filepath.Join(mp, "app", "httpinfo.go"):                     "app.httpinfo.go.tmpl",
		filepath.Join(mp, "app", "presethttpinfo.go"):               "app.presethttpinfo.go.tmpl",
		filepath.Join("config", "httpinfo.toml"):                    "httpinfo.toml.tmpl",
		filepath.Join("system", "constants", "presethttpinfo.toml"): "presethttpinfo.toml.tmpl",
	}
	return task.RenderFiles(filesToRender, nil)
}

var HTTPInfoModule = &HTTPInfo{}

func init() {
	app.Register(HTTPInfoModule)
}
