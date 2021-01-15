package notificationpublisher

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/project"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Queue struct {
	app.BasicModule
	SlienceMode bool
}

type renderData struct {
	Name           *name.Name
	InstallSession bool
}

func (m *Queue) ID() string {
	return "github.com/herb-go/herb-go/modules/notificationpublisher"
}

func (m *Queue) Cmd() string {
	return "notificationpublisher"
}

func (m *Queue) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s notificationpublisher <name>.
Create notification queue module and config files.
Default name is "notificationpublisher".
File below will be created:
	config/<name>.toml
	system/configskeleton/<name>.toml
	src/vendor/modules/app/<name>.go
	src/vendor/modules/<name>/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Queue) Desc(a *app.Application) string {
	return "Create notificationpublisher module and config files"
}
func (m *Queue) Group(a *app.Application) string {
	return "Notification"
}
func (m *Queue) Init(a *app.Application, args *[]string) error {
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

func (m *Queue) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name

	if len(args) == 0 {
		fmt.Println("No notification queue module name given.\"notificationpublisher\" is used")
		n, err = name.New(true, "notificationpublisher")
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
	task := tools.NewTask(filepath.Join(app, "/modules/notificationpublisher/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Queue  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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
	err = tools.CopyIfNotExist(filepath.Join(task.SrcFolder, "drivers.notificationqueue.go"), mp, "drivers", "notificationqueue.go")
	if err != nil {
		return err
	}

	return task.Exec()

}

func (m *Queue) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {
	filesToRender := map[string]string{
		filepath.Join(mp, n.LowerPath("init.go")):                                        "notificationpublisher.modules.go.tmpl",
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):                   "notificationpublisher.toml.tmpl",
		filepath.Join("system", "configskeleton", n.LowerWithParentDotSeparated+".toml"): "notificationpublisher.toml.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"):                    "app.notificationpublisher.go.tmpl",
	}
	data := renderData{
		Name: n,
	}
	return task.RenderFiles(filesToRender, data)
}

var QueueModule = &Queue{}

func init() {
	app.Register(QueueModule)
}
