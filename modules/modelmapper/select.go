package modelmapper

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/herb-go/datasource/sql/db"
	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/util/cli/name"
	"github.com/herb-go/util/config/tomlconfig"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type SelectJoined struct {
	Name    *name.Name
	Prefix  string
	Columns *ModelColumns
}

type Select struct {
	Database string
	Location string
	app.BasicModule
	QueryID     string
	SlienceMode bool
	Prefix      string
	joined      string
	WithRead    bool
	WithList    bool
	WithPager   bool
	User        string
	WithUser    bool
	UserModule  *name.Name
	Joined      []*SelectJoined
}

func (m *Select) ID() string {
	return "github.com/herb-go/herb-go/modules/modelmapper/select"
}

func (m *Select) Cmd() string {
	return "modelmapperselect"
}

func (m *Select) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s modelmapperselect [name].
Create model select code.
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Select) Desc(a *app.Application) string {
	return "Create model select code."
}
func (m *Select) Group(a *app.Application) string {
	return "Model"
}

func (m *Select) GetColumn(table string, prefix string) (*ModelColumns, error) {
	conn := db.New()
	c := db.Config{}
	tomlconfig.MustLoad(source.File("./config/"+m.Database+".toml"), &c)
	conn.SetDriver(c.Driver)
	err := c.ApplyTo(conn)
	if err != nil {
		return nil, err
	}
	return NewModelColumns(conn, m.Database, table, m.Prefix)
}

func (m *Select) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().StringVar(&m.Database, "database", "database",
		`database module name. 
	`)
	m.FlagSet().StringVar(&m.Location, "location", "modelmappers",
		`default model code location. 
	`)
	m.FlagSet().StringVar(&m.User, "user", "",
		`create form with given user module. 
	`)
	m.FlagSet().StringVar(&m.QueryID, "id", "",
		`moder mapper select id. 
	`)
	m.FlagSet().StringVar(&m.Prefix, "prefix", "",
		`table field prefix. 
	`)
	m.FlagSet().StringVar(&m.joined, "join", "",
		`joined models.Format: [location]/<tablename>|[prefix],[location]/<tablename>#[prefix]
	`)
	m.FlagSet().BoolVar(&m.WithRead, "withread", false, "Whether create model read code")
	m.FlagSet().BoolVar(&m.WithList, "withlist", false, "Whether create model list code")
	m.FlagSet().BoolVar(&m.WithPager, "withpager", false, "Whether create model pager code")

	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Select) ParseJoined(mp string) error {
	var err error
	joined := strings.TrimSpace(m.joined)
	if joined == "" {
		return nil
	}
	joinedlist := strings.Split(joined, ",")
	m.Joined = []*SelectJoined{}
	for _, v := range joinedlist {
		j := &SelectJoined{}
		v = strings.TrimSpace(v)
		list := strings.SplitN(v, "#", 2)
		j.Name, err = name.New(true, v)
		if err != nil {
			return err
		}
		if j.Name.Parents == "" && m.Location != "" {
			j.Name, err = name.New(true, m.Location+"/"+j.Name.Raw)
			if err != nil {
				return err
			}
		}

		result, err := tools.FileExists(mp, j.Name.LowerWithParentPath, "models", j.Name.Lower+".go")
		if err != nil {
			return err
		}
		if !result {
			return fmt.Errorf("Model file \"%s\"not found", filepath.Join(mp, j.Name.LowerWithParentPath, "models", j.Name.Lower+".go"))
		}
		if len(list) > 1 {
			j.Prefix = list[1]
		}
		j.Columns, err = m.GetColumn(j.Name.Raw, j.Prefix)
		if err != nil {
			return err
		}

		m.Joined = append(m.Joined, j)
	}
	return nil
}

func (m *Select) Question(a *app.Application) error {
	var err error
	if m.SlienceMode {
		return nil
	}
	err = QuestionWithRead.ExecIf(a, !m.WithRead, &m.WithRead)
	if err != nil {
		return err
	}
	err = QuestionWithList.ExecIf(a, !m.WithList, &m.WithList)
	if err != nil {
		return err
	}
	if m.WithList {
		err = QuestionWithPager.ExecIf(a, !m.WithPager, &m.WithPager)
		if err != nil {
			return err
		}
	}

	return nil
}
func (m *Select) Exec(a *app.Application, args []string) error {
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
	if m.QueryID == "" {
		m.QueryID = "select"
	}
	qn, err := name.New(false, m.QueryID)
	if err != nil {
		return err
	}

	mp, err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}

	if m.User != "" {
		m.WithUser = true
		m.UserModule, err = name.New(true, m.User)
		if err != nil {
			return err
		}
		result, err := tools.FileExists(mp, m.UserModule.LowerWithParentPath, "init.go")
		if err != nil {
			return err
		}
		if !result {
			return fmt.Errorf("User module code \"%s\"not found", filepath.Join(mp, m.UserModule.LowerWithParentPath, "init.go"))
		}
	}

	err = m.ParseJoined(mp)
	if err != nil {
		return err
	}

	file := filepath.Join(mp, n.LowerPath("models"), n.Lower+".go")
	result, err := tools.FileExists(file)
	if err != nil {
		return err
	}
	if !result {
		return fmt.Errorf("model file \"%s\"not found", file)
	}
	mc, err := m.GetColumn(n.Raw, m.Prefix)
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
		a.Printf("ModelMapper  select \"%s/%s\" created.\n", n.LowerWithParentDotSeparated, qn.Lower)
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

func (m *Select) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name, id *name.Name, mc *ModelColumns) error {
	modelmodule := n.LowerWithParentPath
	filesToRender := map[string]string{}
	filesToRender[filepath.Join(mp, n.LowerPath(id.Lower, "forms"), n.Lower+"form.go")] = "selectform.go.tmpl"
	filesToRender[filepath.Join(mp, n.LowerPath(id.Lower, "actions"), n.Lower+"action.go")] = "selectaction.go.tmpl"
	filesToRender[filepath.Join(mp, n.LowerPath(id.Lower, "viewmodels"), n.Lower+"viewmodel.go")] = "selectviewmodel.go.tmpl"
	filesToRender[filepath.Join(mp, n.LowerPath(id.Lower, "results"), n.Lower+"result.go")] = "selectresult.go.tmpl"
	data := map[string]interface{}{
		"Name":      n,
		"Columns":   mc,
		"ID":        id,
		"Module":    modelmodule,
		"Joined":    m.Joined,
		"HasJoined": len(m.Joined) > 0,
		"Confirmed": m,
	}
	return task.RenderFiles(filesToRender, data)
}

var SelectModule = &Select{}
