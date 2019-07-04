package module

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/config"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Module struct {
	app.BasicModule
	Level       string
	AutoConfirm bool
}

func (m *Module) ID() string {
	return "github.com/herb-go/herb-go/modules/module"
}

func (m *Module) Cmd() string {
	return "module"
}

func (m *Module) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s module [name].
Create module.
File below will be created:
	src/vendor/modules/<name>/init.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Module) Desc(a *app.Application) string {
	return "Create new module code"
}

func (m *Module) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.AutoConfirm, "y", false, "Whether auto confirm")
	m.FlagSet().StringVar(&m.Level, "level", "900",
		`Module prefix.All modules are sorted by prefix when loading.
		`)

	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Module) Question(a *app.Application) error {
	return nil
}
func (m *Module) Exec(a *app.Application, args []string) error {
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
	err = config.ErrorIfNotInAppFolder(a.Cwd)
	if err != nil {
		return err
	}
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}

	task := tools.NewTask(filepath.Join(app, "/modules/module/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Module  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *Module) Render(a *app.Application, appPath string, task *tools.Task, n *name.Name) error {
	return task.Render("module.go.tmpl", filepath.Join("src", "vendor", "modules", n.LowerPath("init.go")), map[string]interface{}{"Name": n, "Level": m.Level})
}

var ModuelModule = &Module{}

func init() {
	app.Register(ModuelModule)
}
