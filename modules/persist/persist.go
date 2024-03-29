package persist

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/overseers"
	"github.com/herb-go/herb-go/modules/project"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Persist struct {
	app.BasicModule
	SlienceMode bool
}

type renderData struct {
	Name           *name.Name
	InstallSession bool
}

func (m *Persist) ID() string {
	return "github.com/herb-go/herb-go/modules/persist"
}

func (m *Persist) Cmd() string {
	return "persist"
}

func (m *Persist) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s persist <name>.
Create persist module and config files.
Default name is "persist".
File below will be created:
	config/<name>.toml
	system/configskeleton/<name>.toml
	src/modules/app/<name>.go
	src/modules/<name>/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Persist) Desc(a *app.Application) string {
	return "Create persist module and config files"
}
func (m *Persist) Group(a *app.Application) string {
	return "Data"
}
func (m *Persist) Init(a *app.Application, args *[]string) error {
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

func (m *Persist) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name

	if len(args) == 0 {
		fmt.Println("No persist module name given.\"persistdata\" is used")
		n, err = name.New(true, "persistdata")
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
	err = overseers.OverseerModule.Exec(a, []string{"-s", "persist"})
	if err != nil {
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/persist/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Persist  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *Persist) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {
	err := tools.CopyIfNotExist(filepath.Join(task.SrcFolder, "persist.hired.go.tmpl"), mp, "hired", "persist.go")
	if err != nil {
		return err
	}
	filesToRender := map[string]string{
		filepath.Join(mp, n.LowerPath("init.go")):                                        "persist.modules.go.tmpl",
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):                   "persist.toml.tmpl",
		filepath.Join("system", "configskeleton", n.LowerWithParentDotSeparated+".toml"): "persist.toml.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"):                    "app.persist.go.tmpl",
	}
	data := renderData{
		Name: n,
	}
	return task.RenderFiles(filesToRender, data)
}

var PersistModule = &Persist{}

func init() {
	app.Register(PersistModule)
}
