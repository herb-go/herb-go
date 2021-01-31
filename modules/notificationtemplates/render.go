package notificationtemplates

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/driver"
	"github.com/herb-go/herb-go/modules/project"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Render struct {
	app.BasicModule
	SlienceMode bool
}

type renderData struct {
	Name           *name.Name
	InstallSession bool
}

func (m *Render) ID() string {
	return "github.com/herb-go/herb-go/modules/notificationtemplates"
}

func (m *Render) Cmd() string {
	return "notificationtemplates"
}

func (m *Render) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s notificationtemplates <name>.
Create notificationtemplates module and config files.
Default name is "notificationtemplates".
File below will be created:
	config/<name>.toml
	system/configskeleton/<name>.toml
	src/vendor/modules/app/<name>.go
	src/vendor/modules/<name>/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Render) Desc(a *app.Application) string {
	return "Create notificationtemplates module and config files"
}
func (m *Render) Group(a *app.Application) string {
	return "Notification"
}
func (m *Render) Init(a *app.Application, args *[]string) error {
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

func (m *Render) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name

	if len(args) == 0 {
		fmt.Println("No notificationtemplates module name given.\"notificationtemplates\" is used")
		n, err = name.New(true, "notificationtemplates")
	} else {
		n, err = name.New(true, args...)
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
	task := tools.NewTask(filepath.Join(app, "/modules/notificationtemplates/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("notificationtemplates  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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
	err = driver.DriverModule.Exec(a, []string{"-s", "texttemplate"})
	if err != nil {
		return err
	}

	return task.Exec()

}

func (m *Render) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {
	filesToRender := map[string]string{
		filepath.Join(mp, n.LowerPath("init.go")):                     "notificationtemplates.modules.go.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"): "app.notificationtemplates.go.tmpl",
	}
	data := renderData{
		Name: n,
	}
	err := task.CopyFiles(map[string]string{
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):                   "notificationtemplates.toml",
		filepath.Join("system", "configskeleton", n.LowerWithParentDotSeparated+".toml"): "notificationtemplates.toml",
	})
	if err != nil {
		return err
	}
	return task.RenderFiles(filesToRender, data)
}

var notificationtemplatesModule = &Render{}

func init() {
	app.Register(notificationtemplatesModule)
}
