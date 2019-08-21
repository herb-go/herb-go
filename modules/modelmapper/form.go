package modelmapper

import (
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

type Form struct {
	Database string
	Location string
	app.BasicModule
	QueryID     string
	SlienceMode bool
	Prefix      string
	Member      string
	WithMember  bool
	MemberName  *name.Name
}

func (m *Form) ID() string {
	return "github.com/herb-go/herb-go/modules/modelmapper/form"
}

func (m *Form) Cmd() string {
	return "modelmapperform"
}

func (m *Form) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s modelmapperform [name].
Create model form code.
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Form) Desc(a *app.Application) string {
	return "Create model form code."
}
func (m *Form) Group(a *app.Application) string {
	return "Model"
}

func (m *Form) GetColumn(table string) (*ModelColumns, error) {
	conn := db.New()
	c := db.Config{}
	tomlconfig.MustLoad(util.File("./config/"+m.Database+".toml"), &c)
	conn.SetDriver(c.Driver)
	err := c.ApplyTo(conn)
	if err != nil {
		return nil, err
	}
	return NewModelColumns(conn, m.Database, table, m.Prefix)
}

func (m *Form) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().StringVar(&m.Database, "database", "database",
		`database module name. 
	`)
	m.FlagSet().StringVar(&m.Location, "location", "modelmappers",
		`default model code location. 
	`)
	m.FlagSet().StringVar(&m.Member, "member", "",
		`create form with given member. 
	`)
	m.FlagSet().StringVar(&m.QueryID, "id", "",
		`moder mapper form id. 
	`)
	m.FlagSet().StringVar(&m.Prefix, "prefix", "",
		`table field prefix. 
	`)
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}

func (m *Form) Question(a *app.Application) error {
	return nil
}
func (m *Form) Exec(a *app.Application, args []string) error {
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
	if n.Parents == "" && m.Location != "" {
		n, err = name.New(true, m.Location+"/"+n.Raw)
		if err != nil {
			return err
		}
	}
	qn, err := name.New(false, m.QueryID)
	if err != nil {
		return err
	}

	mp, err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}

	if m.Member != "" {
		m.WithMember = true
		m.MemberName, err = name.New(true, m.Member)
		if err != nil {
			return err
		}
		result, err := tools.FileExists(mp, m.MemberName.LowerWithParentPath, "init.go")
		if err != nil {
			return err
		}
		if !result {
			return fmt.Errorf("Member file \"%s\"not found", filepath.Join(mp, m.MemberName.LowerWithParentPath, "init.go"))
		}
	}

	file := filepath.Join(mp, n.LowerPath("models"), n.Lower+".go")
	result, err := tools.FileExists(file)
	if err != nil {
		return err
	}
	if !result {
		return fmt.Errorf("model file \"%s\"not found", file)
	}
	mc, err := m.GetColumn(n.Raw)
	if err != nil {
		return err
	}
	if !mc.IsSinglePrimayKey() || (mc.PrimaryKeys[0].ColumnType != "string" && mc.PrimaryKeys[0].ColumnType != "int" && mc.PrimaryKeys[0].ColumnType != "int64") {
		return ErrUnsupportedPK
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

	err = m.Render(a, a.Cwd, mp, task, n, qn, mc)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("ModelMapper  form \"%s %s\" created.\n", n.LowerWithParentDotSeparated, qn.Lower)
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

func (m *Form) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name, id *name.Name, mc *ModelColumns) error {
	modelmodule := n.LowerWithParentPath
	filesToRender := map[string]string{}
	filesToRender[filepath.Join(mp, n.LowerPath(id.Lower, "forms"), n.Lower+"form.go")] = "form.go.tmpl"
	filesToRender[filepath.Join(mp, n.LowerPath(id.Lower, "actions"), n.Lower+"action.go")] = "action.go.tmpl"
	data := map[string]interface{}{
		"Name":      n,
		"Columns":   mc,
		"ID":        id,
		"Module":    modelmodule,
		"Confirmed": m,
	}
	return task.RenderFiles(filesToRender, data)
}

var FormModule = &Form{}
