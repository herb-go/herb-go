package modelmapper

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/util"
	"github.com/herb-go/util/cli/name"
	"github.com/herb-go/util/config/tomlconfig"

	"github.com/herb-go/herb-go/modules/module"
	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

var QuestionCRUD = tools.NewTrueOrFalseQuestion("Do you want to install standard \"CRUD\" components")
var QuestionWithCreate = tools.NewTrueOrFalseQuestion("Do you want to create model \"Create\" component?")
var QuestionWithRead = tools.NewTrueOrFalseQuestion("Do you want to create model \"Read\" component?")
var QuestionWithUpdate = tools.NewTrueOrFalseQuestion("Do you want to create model \"Update\" component?")
var QuestionWithDelete = tools.NewTrueOrFalseQuestion("Do you want to create model \"Delete\" component?")
var QuestionWithList = tools.NewTrueOrFalseQuestion("Do you want to create model \"List\" component?")
var QuestionWithPager = tools.NewTrueOrFalseQuestion("Do you want to use pager for  \"List\" component?")
var QuestionCreateForm = tools.NewTrueOrFalseQuestion("Do you want to create model forms?")
var QuestionCreateAction = tools.NewTrueOrFalseQuestion("Do you want to create model actions?")
var QuestionCreateOutput = tools.NewTrueOrFalseQuestion("Do you want to create model output class")

type ModelMapper struct {
	app.BasicModule
	Database     string
	CreateForm   bool
	CreateOutput bool
	CreateAction bool
	WithCreate   bool
	WithRead     bool
	WithUpdate   bool
	WithDelete   bool
	WithList     bool
	WithPager    bool
	SlienceMode  bool
}

func (m *ModelMapper) ID() string {
	return "github.com/herb-go/herb-go/modules/modelmapper"
}

func (m *ModelMapper) Cmd() string {
	return "modelmapper"
}

func (m *ModelMapper) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s modelmapper <name>.
Create model module and config files.
File below will be created:
	src/vendor/modules/<name>/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *ModelMapper) Desc(a *app.Application) string {
	return "Create model mapper and config files."
}
func (m *ModelMapper) Group(a *app.Application) string {
	return "Model"
}
func (m *ModelMapper) GetColumn(table string) (*ModelColumns, error) {
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
func (m *ModelMapper) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.SlienceMode, "s", false, "Slience mode")
	m.FlagSet().StringVar(&m.Database, "database", "database",
		`database module name. 
	`)

	crud := m.FlagSet().Bool("crud", false, "Whether create all CRUD codes")
	m.FlagSet().BoolVar(&m.CreateAction, "createaction", false, "Whether create model actions")
	m.FlagSet().BoolVar(&m.CreateForm, "createform", false, "Whether create model forms")
	m.FlagSet().BoolVar(&m.CreateOutput, "createoutput", false, "Whether create model output class")
	m.FlagSet().BoolVar(&m.WithCreate, "withcreate", false, "Whether create model create code")
	m.FlagSet().BoolVar(&m.WithRead, "withread", false, "Whether create model read code")
	m.FlagSet().BoolVar(&m.WithUpdate, "withupdate", false, "Whether create model update code")
	m.FlagSet().BoolVar(&m.WithDelete, "withdelete", false, "Whether create model delete code")
	m.FlagSet().BoolVar(&m.WithList, "withlist", false, "Whether create model list code")
	m.FlagSet().BoolVar(&m.WithPager, "withpager", false, "Whether create model pager code")

	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	if *crud {
		m.WithCreate = true
		m.WithRead = true
		m.WithUpdate = true
		m.WithDelete = true
		m.CreateAction = true
		m.CreateForm = true
		m.CreateOutput = true
		m.WithList = true
		m.WithPager = true
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *ModelMapper) Question(a *app.Application, mc *ModelColumns) error {
	if m.SlienceMode {
		return nil
	}
	if len(mc.PrimaryKeys) == 1 && (mc.PrimaryKeys[0].ColumnType == "string" || mc.PrimaryKeys[0].ColumnType == "int") {
		crud := m.WithCreate && m.WithRead && m.WithUpdate && m.WithDelete
		err := QuestionCRUD.ExecIf(a, !crud, &crud)
		if err != nil {
			return err
		}
		if crud {
			m.WithCreate = true
			m.WithRead = true
			m.WithUpdate = true
			m.WithDelete = true
			m.WithList = true
			m.WithPager = true
			m.CreateAction = true
			m.CreateForm = true
			m.CreateOutput = true
		}
		err = QuestionWithCreate.ExecIf(a, mc.CanCreate() && !m.WithCreate, &m.WithCreate)
		if err != nil {
			return err
		}
		if mc.HasPrimayKey() {
			err = QuestionWithRead.ExecIf(a, !m.WithRead, &m.WithRead)
			if err != nil {
				return err
			}
			err = QuestionWithUpdate.ExecIf(a, !m.WithUpdate, &m.WithUpdate)
			if err != nil {
				return err
			}
			err = QuestionWithDelete.ExecIf(a, !m.WithDelete, &m.WithDelete)
			if err != nil {
				return err
			}
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
		if m.WithCreate || m.WithRead || m.WithUpdate || m.WithDelete || m.WithList {
			err = QuestionCreateForm.ExecIf(a, !m.CreateForm, &m.CreateForm)
			if err != nil {
				return err
			}
			err = QuestionCreateAction.ExecIf(a, !m.CreateAction, &m.CreateAction)
			if err != nil {
				return err
			}
		}
		err = QuestionCreateOutput.ExecIf(a, !m.CreateOutput, &m.CreateOutput)
		if err != nil {
			return err
		}
	}
	return nil
}
func (m *ModelMapper) Exec(a *app.Application, args []string) error {
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
	err = m.Question(a, mc)
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
		a.Printf("ModelMapper  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *ModelMapper) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name, mc *ModelColumns) error {
	modelmodule := n.LowerWithParentPath
	modulepath := filepath.Join(mp, modelmodule)
	exists, err := tools.FileExists(modulepath)
	if err != nil {
		return err
	}
	if !exists {
		err = module.ModuleModule.Exec(a, []string{modelmodule})
		if err != nil {
			return err
		}
	}
	filesToRender := map[string]string{
		filepath.Join(mp, n.LowerPath("models"), n.Lower+"queries.go"): "modelqueries.go.tmpl",
		filepath.Join(mp, n.LowerPath("models"), n.Lower+"columns.go"): "modelcolumns.go.tmpl",
		filepath.Join(mp, n.LowerPath("models"), n.Lower+"fields.go"):  "modelfields.go.tmpl",
		filepath.Join(mp, n.LowerPath("models"), n.Lower+".go"):        "model.go.tmpl",
	}
	if len(mc.PrimaryKeys) == 1 && (mc.PrimaryKeys[0].ColumnType == "string" || mc.PrimaryKeys[0].ColumnType == "int") {
		if m.CreateForm {
			filesToRender[filepath.Join(mp, n.LowerPath("forms"), n.Lower+"form.go")] = "modelform.go.tmpl"
		}
		if m.CreateAction {
			filesToRender[filepath.Join(mp, n.LowerPath("actions"), n.Lower+"action.go")] = "modelaction.go.tmpl"
		}
		if m.CreateOutput {
			filesToRender[filepath.Join(mp, n.LowerPath("outputs"), n.Lower+"output.go")] = "modeloutput.go.tmpl"
		}
	}
	data := map[string]interface{}{
		"Name":      n,
		"Columns":   mc,
		"Module":    modelmodule,
		"Confirmed": m,
	}
	return task.RenderFiles(filesToRender, data)
}

var ModelMapperModule = &ModelMapper{}

func init() {
	app.Register(ModelMapperModule)
}
