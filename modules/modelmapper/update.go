package modelmapper

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/util"
	"github.com/herb-go/util/cli/name"
	"github.com/herb-go/util/config/tomlconfig"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

var QuestionUpdateColumns = tools.NewTrueOrFalseQuestion("Do you want to update model columns code?")
var QuestionConfirmBackupColumns = tools.NewTrueOrFalseQuestion("Have you backed up model columns code?")
var QuestionUpdateFields = tools.NewTrueOrFalseQuestion("Do you want to update model columns code?")
var QuestionConfirmBackupFields = tools.NewTrueOrFalseQuestion("Have you backed up model columns code?")
var QuestionUpdateQueries = tools.NewTrueOrFalseQuestion("Do you want to update model queries code?")
var QuestionConfirmBackupQueries = tools.NewTrueOrFalseQuestion("Have you backed up model queries code?")

type Update struct {
	Database string
	app.BasicModule
	UpdateColumns bool
	UpdateFields  bool
	UpdateQueries bool
}

func (m *Update) ID() string {
	return "github.com/herb-go/herb-go/modules/modelmapper/update"
}

func (m *Update) Cmd() string {
	return "modelmapperupdate"
}

func (m *Update) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s modelmapperupdate [name].
Update model module and config files.
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Update) Desc(a *app.Application) string {
	return "Update model module and config files."
}
func (m *Update) Group(a *app.Application) string {
	return "Model"
}

func (m *Update) GetColumn(table string) (*ModelColumns, error) {
	conn := db.New()
	c := db.Config{}
	tomlconfig.MustLoad(util.File("./config/"+m.Database+".toml"), &c)
	conn.SetDriver(c.Driver)
	err := c.ApplyTo(conn)
	if err != nil {
		return nil, err
	}
	return NewModelCulumns(conn, m.Database, table)
}

func (m *Update) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().StringVar(&m.Database, "database", "database",
		`database module name. 
	`)
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}

func (m *Update) confirmFile(a *app.Application, q *tools.Question, qbackuped *tools.Question, result *bool) error {
	var confirmed bool
	err := q.ExecIf(a, true, result)
	if err != nil {
		return err
	}
	if *result {
		err = qbackuped.ExecIf(a, true, &confirmed)
		if err != nil {
			return err
		}
		if !confirmed {
			return errors.New("You shoud backup file firest")
		}
	}
	return nil
}
func (m *Update) Question(a *app.Application) error {
	var err error
	err = m.confirmFile(a, QuestionUpdateColumns, QuestionConfirmBackupColumns, &m.UpdateColumns)
	if err != nil {
		return err
	}
	err = m.confirmFile(a, QuestionUpdateFields, QuestionConfirmBackupFields, &m.UpdateFields)
	if err != nil {
		return err
	}
	err = m.confirmFile(a, QuestionUpdateQueries, QuestionConfirmBackupQueries, &m.UpdateQueries)
	if err != nil {
		return err
	}

	return nil
}
func (m *Update) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		a.PrintModuleHelp(m)
		return nil
	}
	n, err := name.New(true, args...)
	if err != nil {
		return err
	}

	mp, err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}
	mc, err := m.GetColumn(n.Raw)
	if err != nil {
		return err
	}
	err = m.Question(a)
	if err != nil {
		return err
	}
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}

	task := tools.NewTask(filepath.Join(app, "/modules/modelmapper/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n, mc)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("ModelMapper  \"%s\" updated.\n", n.LowerWithParentDotSeparated)
		return nil
	})
	err = task.ErrosIfAnyFileNotExists()
	if err != nil {
		return err
	}
	ok, err := task.ConfirmIf(a, true)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return task.Exec()
}

func (m *Update) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name, mc *ModelColumns) error {
	modelmodule := n.LowerWithParentPath
	filesToRender := map[string]string{}
	if m.UpdateColumns {
		filesToRender[filepath.Join(mp, n.LowerPath("models"), n.Lower+"columns.go")] = "modelcolumns.go.tmpl"
	}
	if m.UpdateFields {
		filesToRender[filepath.Join(mp, n.LowerPath("models"), n.Lower+"fields.go")] = "modelfields.go.tmpl"
	}
	if m.UpdateQueries {
		filesToRender[filepath.Join(mp, n.LowerPath("models"), n.Lower+"queries.go")] = "modelqueries.go.tmpl"
	}

	data := map[string]interface{}{
		"Name":      n,
		"Columns":   mc,
		"Module":    modelmodule,
		"Confirmed": m,
	}
	return task.RenderFiles(filesToRender, data)
}

var UpdateModule = &Update{}

func init() {
}
