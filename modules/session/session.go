package session


import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Session struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Session) ID() string {
	return "github.com/herb-go/herb-go/modules/session"
}

func (m *Session) Cmd() string {
	return "session"
}

func (m *Session) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s session <name>.
Create session module and config files.
Default name is "session".
File below will be created:
	config/<name>.toml
	system/config.examples/<name>.toml
	src/vendor/modules/app/<name>.go
	src/vendor/modules/<name>/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Session) Desc(a *app.Application) string {
	return "Create session module and config files."
}
func (m *Session) Group(a *app.Application) string {
	return "Auth"
}
func (m *Session) Init(a *app.Application, args *[]string) error {
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
func (m *Session) Question(a *app.Application) error {
	return nil
}
func (m *Session) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name

	if len(args) == 0 {
		fmt.Println("No session module name given.\"session\" is used")
		n ,err= name.New(true,"session")
	}else{
		n, err = name.New(true, args...)
	}
	if err != nil {
		return err
	}

	mp,err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}

	task := tools.NewTask(filepath.Join(app, "/modules/session/resources"), a.Cwd)

	err = m.Render(a, a.Cwd,mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Session  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *Session) Render(a *app.Application, appPath string,mp string, task *tools.Task, n *name.Name) error {
	filesToRender := map[string]string{
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):"session.toml.tmpl",
		filepath.Join("system", "config.examples", n.LowerWithParentDotSeparated+".toml"):"session.toml.tmpl",
		filepath.Join(mp,n.LowerPath(n.Lower+".go")):           "session.modules.go.tmpl",
		filepath.Join(mp,"app", n.LowerWithParentDotSeparated+".go"):           "session.go.tmpl",
	}
	return task.RenderFiles(filesToRender, n)
}

var SessionModule = &Session{}

func init() {
	app.Register(SessionModule)
}
