package config

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Translate struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Translate) ID() string {
	return "github.com/herb-go/herb-go/modules/translate"
}

func (m *Translate) Cmd() string {
	return "translate"
}

func (m *Translate) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s translate.
Create translate config and code.
File below will be created:
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Translate) Desc(a *app.Application) string {
	return "Create translate config and code"
}
func (m *Translate) Group(a *app.Application) string {
	return "Translate"
}
func (m *Translate) Init(a *app.Application, args *[]string) error {
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
func (m *Translate) Question(a *app.Application) error {
	return nil
}
func (m *Translate) Exec(a *app.Application, args []string) error {
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

	task := tools.NewTask(filepath.Join(app, "/modules/translate/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Translate module  created.\n")
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

func (m *Translate) Render(a *app.Application, appPath string, mp string, task *tools.Task) error {

	filesToRender := map[string]string{
		filepath.Join(mp, "app", "translate.go"):              "translate.go.tmpl",
		filepath.Join(mp, "translate", "init.go"):             "translate/translate.module.go.tmpl",
		filepath.Join(mp, "test", "translate.go"):             "translate.test.go.tmpl",
		filepath.Join("system", "messages", "translate.toml"): "messages/translate.toml",
		filepath.Join("system", "messages", "readme.md"):      "messages/readme.md",
		filepath.Join("src", "translate.go"):                  "translate.src.go.tmpl",
	}
	return task.RenderFiles(filesToRender, nil)
}

var TranslateModule = &Translate{}

func init() {
	app.Register(TranslateModule)
}
