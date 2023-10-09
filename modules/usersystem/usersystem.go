package usersystem

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/overseers"
	"github.com/herb-go/herb-go/modules/project"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type UserSystem struct {
	app.BasicModule
	SlienceMode bool
}

type renderData struct {
	Name           *name.Name
	InstallSession bool
}

func (m *UserSystem) ID() string {
	return "github.com/herb-go/herb-go/modules/usersystem"
}

func (m *UserSystem) Cmd() string {
	return "user"
}

func (m *UserSystem) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s usersystem <name>.
Create usersystem module and config files.
Default name is "usersystem".
File below will be created:
	config/<name>.toml
	system/configskeleton/<name>.toml
	src/modules/app/<name>.go
	src/modules/<name>/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *UserSystem) Desc(a *app.Application) string {
	return "Create usersystem module and config files"
}
func (m *UserSystem) Group(a *app.Application) string {
	return "Auth"
}
func (m *UserSystem) Init(a *app.Application, args *[]string) error {
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

func (m *UserSystem) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name

	if len(args) == 0 {
		fmt.Println("No usersystem module name given.\"members\" is used")
		n, err = name.New(true, "members")
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
	err = overseers.OverseerModule.Exec(a, []string{"-s", "usersystem"})
	if err != nil {
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/usersystem/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("UserSystem  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *UserSystem) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {
	err := tools.CopyIfNotExist(filepath.Join(task.SrcFolder, "usersystem.hired.go.tmpl"), mp, "hired", "usersystem.go")
	if err != nil {
		return err
	}
	filesToRender := map[string]string{
		filepath.Join(mp, n.LowerPath("init.go")):                                        "usersystem.modules.go.tmpl",
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):                   "usersystem.toml.tmpl",
		filepath.Join("system", "configskeleton", n.LowerWithParentDotSeparated+".toml"): "usersystem.toml.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"):                    "app.usersystem.go.tmpl",
		filepath.Join(mp, n.LowerPath("middleware.go")):                                  "middleware.go.tmpl",
	}
	data := renderData{
		Name: n,
	}
	return task.RenderFiles(filesToRender, data)
}

var UserSystemModule = &UserSystem{}

func init() {
	app.Register(UserSystemModule)
}
