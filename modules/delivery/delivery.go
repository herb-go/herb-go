package delivery

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/project"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Delivery struct {
	app.BasicModule
	SlienceMode bool
}

type renderData struct {
	Name           *name.Name
	InstallSession bool
}

func (m *Delivery) ID() string {
	return "github.com/herb-go/herb-go/modules/delivery"
}

func (m *Delivery) Cmd() string {
	return "delivery"
}

func (m *Delivery) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s delivery <name>.
Create delivery module and config files.
Default name is "delivery".
File below will be created:
	config/<name>.toml
	system/configskeleton/<name>.toml
	src/modules/app/<name>.go
	src/modules/<name>/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Delivery) Desc(a *app.Application) string {
	return "Create delivery module and config files"
}
func (m *Delivery) Group(a *app.Application) string {
	return "Notification"
}
func (m *Delivery) Init(a *app.Application, args *[]string) error {
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

func (m *Delivery) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name

	if len(args) == 0 {
		fmt.Println("No delivery module name given.\"delivery\" is used")
		n, err = name.New(true, "delivery")
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
	task := tools.NewTask(filepath.Join(app, "/modules/delivery/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Delivery  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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
	err = tools.CopyIfNotExist(filepath.Join(task.SrcFolder, "drivers.delivery.go"), mp, "drivers", "delivery.go")
	if err != nil {
		return err
	}

	return task.Exec()

}

func (m *Delivery) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {
	filesToRender := map[string]string{
		filepath.Join(mp, n.LowerPath("init.go")):                                        "delivery.modules.go.tmpl",
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):                   "delivery.toml.tmpl",
		filepath.Join("system", "configskeleton", n.LowerWithParentDotSeparated+".toml"): "delivery.toml.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"):                    "app.delivery.go.tmpl",
	}
	data := renderData{
		Name: n,
	}
	return task.RenderFiles(filesToRender, data)
}

var DeliveryModule = &Delivery{}

func init() {
	app.Register(DeliveryModule)
}
