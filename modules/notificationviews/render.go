package notificationviews

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/driver"
	"github.com/herb-go/herb-go/modules/project"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type View struct {
	app.BasicModule
	SlienceMode bool
}

type renderData struct {
	Name           *name.Name
	InstallSession bool
}

func (m *View) ID() string {
	return "github.com/herb-go/herb-go/modules/notificationviews"
}

func (m *View) Cmd() string {
	return "notificationviews"
}

func (m *View) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s notificationviews <name>.
Create notificationviews module and config files.
Default name is "notificationviews".
File below will be created:
	config/<name>.toml
	system/configskeleton/<name>.toml
	src/modules/app/<name>.go
	src/modules/<name>/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *View) Desc(a *app.Application) string {
	return "Create notificationviews module and config files"
}
func (m *View) Group(a *app.Application) string {
	return "Notification"
}
func (m *View) Init(a *app.Application, args *[]string) error {
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

func (m *View) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name

	if len(args) == 0 {
		fmt.Println("No notificationviews module name given.\"notificationviews\" is used")
		n, err = name.New(true, "notificationviews")
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
	task := tools.NewTask(filepath.Join(app, "/modules/notificationviews/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("notificationviews  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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
	driver.DriverModule.Reset()
	err = driver.DriverModule.Exec(a, []string{"-s", "notificationview"})
	if err != nil {
		return err
	}
	driver.DriverModule.Reset()
	err = driver.DriverModule.Exec(a, []string{"-s", "texttemplate"})
	if err != nil {
		return err
	}
	return task.Exec()

}

func (m *View) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {
	filesToView := map[string]string{
		filepath.Join(mp, n.LowerPath("init.go")):                     "notificationviews.modules.go.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"): "app.notificationviews.go.tmpl",
	}
	data := renderData{
		Name: n,
	}
	err := task.CopyFiles(map[string]string{
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):                   "notificationviews.toml",
		filepath.Join("system", "configskeleton", n.LowerWithParentDotSeparated+".toml"): "notificationviews.toml",
	})
	if err != nil {
		return err
	}
	return task.RenderFiles(filesToView, data)
}

var notificationviewsModule = &View{}

func init() {
	app.Register(notificationviewsModule)
}
