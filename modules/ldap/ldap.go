package database

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Ldap struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Ldap) ID() string {
	return "github.com/herb-go/herb-go/modules/ldap"
}

func (m *Ldap) Cmd() string {
	return "ldap"
}

func (m *Ldap) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s ldap <name>.
Create sql ldap module and config files.
Default name is "ldap".
File below will be created:
	config/<name>.toml
	system/confg.examples/<name>.toml
	src/modules/app/<name>.go
	src/modules/<name>/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Ldap) Desc(a *app.Application) string {
	return "Create ldap module and config files."
}
func (m *Ldap) Group(a *app.Application) string {
	return "Data"
}
func (m *Ldap) Init(a *app.Application, args *[]string) error {
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
func (m *Ldap) Question(a *app.Application) error {
	return nil
}
func (m *Ldap) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name

	if len(args) == 0 {
		fmt.Println("No ldap module name given.\"ldap\" is used")
		n, err = name.New(true, "ldap")
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

	task := tools.NewTask(filepath.Join(app, "/modules/ldap/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Ldap  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *Ldap) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {

	filesToRender := map[string]string{
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):                   "ldap.toml.tmpl",
		filepath.Join("system", "configskeleton", n.LowerWithParentDotSeparated+".toml"): "ldap.toml.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"):                    "ldap.go.tmpl",
	}
	return task.RenderFiles(filesToRender, n)
}

var LdapModule = &Ldap{}

func init() {
	app.Register(LdapModule)
}
