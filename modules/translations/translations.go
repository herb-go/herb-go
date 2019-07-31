package translations

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Translations struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Translations) ID() string {
	return "github.com/herb-go/herb-go/modules/translations"
}

func (m *Translations) Cmd() string {
	return "translations"
}

func (m *Translations) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s translations.
Create translate config and code.
File below will be created:
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Translations) Desc(a *app.Application) string {
	return "Create translations config and code"
}
func (m *Translations) Group(a *app.Application) string {
	return "Translations"
}
func (m *Translations) Init(a *app.Application, args *[]string) error {
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
func (m *Translations) Question(a *app.Application) error {
	return nil
}
func (m *Translations) Exec(a *app.Application, args []string) error {
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

	task := tools.NewTask(filepath.Join(app, "/modules/translations/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Translations module  created.\n")
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

func (m *Translations) Render(a *app.Application, appPath string, mp string, task *tools.Task) error {

	filesToRender := map[string]string{
		filepath.Join(mp, "app", "translations.go"):              "translations.go.tmpl",
		filepath.Join(mp, "translations", "init.go"):             "translations/translations.module.go.tmpl",
		filepath.Join(mp, "test", "translations.go"):             "translations.test.go.tmpl",
		filepath.Join("system", "messages", "translations.toml"): "messages/translations.toml",
		filepath.Join("system", "messages", "readme.md"):         "messages/readme.md",
		filepath.Join("src", "translations.go"):                  "translations.src.go.tmpl",
	}
	return task.RenderFiles(filesToRender, nil)
}

var TranslationsModule = &Translations{}

func init() {
	app.Register(TranslationsModule)
}
