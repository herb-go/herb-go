package database


import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Database struct {
	app.BasicModule
	AutoConfirm bool
}

func (m *Database) ID() string {
	return "github.com/herb-go/herb-go/modules/router"
}

func (m *Database) Cmd() string {
	return "database"
}

func (m *Database) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s databse <name>.
Create sql database module and config files.
Default name is "database".
File below will be created:
	config/<name>.toml
	system/confg.examples/<name>.toml
	src/vendor/modules/app/<name>.go
	src/vendor/modules/<name>/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Database) Desc(a *app.Application) string {
	return "Create sql database module and config files."
}
func (m *Database) Group(a *app.Application) string {
	return "Data"
}
func (m *Database) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.AutoConfirm, "y", false, "Whether auto confirm")
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Database) Question(a *app.Application) error {
	return nil
}
func (m *Database) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name

	if len(args) == 0 {
		fmt.Println("No database module name given.\"database\" is used")
		n ,err= name.New(true,"database")
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

	task := tools.NewTask(filepath.Join(app, "/modules/database/resources"), a.Cwd)

	err = m.Render(a, a.Cwd,mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Database  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *Database) Render(a *app.Application, appPath string,mp string, task *tools.Task, n *name.Name) error {

	filesToRender := map[string]string{
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):"database.toml.tmpl",
		filepath.Join("system", "config.examples", n.LowerWithParentDotSeparated+".toml"):"database.toml.tmpl",
		filepath.Join(mp,n.LowerPath(n.Lower+".go")):           "database.modules.go.tmpl",
		filepath.Join(mp,"app", n.LowerWithParentDotSeparated+".go"):           "database.go.tmpl",
	}
	return task.RenderFiles(filesToRender, n)
}

var DatabaseModule = &Database{}

func init() {
	app.Register(DatabaseModule)
}
