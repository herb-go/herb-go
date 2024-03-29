package cache

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/driver"
	"github.com/herb-go/herb-go/modules/overseers"
	"github.com/herb-go/herb-go/modules/project"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Cache struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Cache) ID() string {
	return "github.com/herb-go/herb-go/modules/cache"
}

func (m *Cache) Cmd() string {
	return "cache"
}

func (m *Cache) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s cache <name>.
Create cache module and config files.
Default name is "cache".
File below will be created:
	config/<name>.toml
	system/configskeleton/<name>.toml
	src/modules/app/<name>.go
	src/modules/<name>/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Cache) Desc(a *app.Application) string {
	return "Create new cache code and config"
}

func (m *Cache) Group(a *app.Application) string {
	return "Data"
}

func (m *Cache) Init(a *app.Application, args *[]string) error {
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
func (m *Cache) Question(a *app.Application) error {
	return nil
}
func (m *Cache) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name
	if len(args) == 0 {
		a.Println("No cache module name given.\"cache\" is used")
		n, err = name.New(true, "cache")
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
	err = overseers.OverseerModule.Exec(a, []string{"-s", "cache"})
	if err != nil {
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/cache/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Cache  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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
	err = task.Exec()
	if err != nil {
		return err
	}
	driver.DriverModule.Reset()
	err = driver.DriverModule.Exec(a, []string{"-s", "kvdb"})
	if err != nil {
		return err
	}
	return nil
}

func (m *Cache) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {

	filesToRender := map[string]string{
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):                   "cache.toml.tmpl",
		filepath.Join("system", "configskeleton", n.LowerWithParentDotSeparated+".toml"): "cache.toml.tmpl",
		filepath.Join(mp, n.LowerPath("init.go")):                                        "cache.modules.go.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"):                    "cache.go.tmpl",
	}
	return task.RenderFiles(filesToRender, n)
}

var Module = &Cache{}

func init() {
	app.Register(Module)
}
