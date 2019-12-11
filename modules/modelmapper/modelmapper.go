package modelmapper

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/util/cli/name"
	"github.com/herb-go/util/config/tomlconfig"

	"github.com/herb-go/herb-go/modules/module"
	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/herb-go/modules/uniqueid"
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
var QuestionCreateViewModel = tools.NewTrueOrFalseQuestion("Do you want to create model view model class?")
var QuestionAutoPk = tools.NewTrueOrFalseQuestion("Do you want to auto fill model primary key with unique id ?")

func NewQuestionCreatedField(field *name.Name) *tools.Question {
	return tools.NewTrueOrFalseQuestion(fmt.Sprintf("Do you want to auto fill model field \"%s\" as created time with time.Unix() ?", field.Raw))
}
func NewQuestionUpdatedField(field *name.Name) *tools.Question {
	return tools.NewTrueOrFalseQuestion(fmt.Sprintf("Do you want to auto fill model field \"%s\" as updated time with time.Unix() ?", field.Raw))
}

type ModelMapper struct {
	app.BasicModule
	Database              string
	Location              string
	QueryID               string
	CreateForm            bool
	CreateViewModel       bool
	CreateAction          bool
	WithCreate            bool
	WithRead              bool
	WithUpdate            bool
	WithDelete            bool
	WithList              bool
	WithPager             bool
	SlienceMode           bool
	Prefix                string
	User                  string
	WithUser              bool
	InstallUniqueID       bool
	AutoPK                bool
	CreatedTimestampField string
	UpdatedTimestampField string
	UserModule            *name.Name
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
	tomlconfig.MustLoad(source.File("./config/"+m.Database+".toml"), &c)
	conn.SetDriver(c.Driver)
	err := c.ApplyTo(conn)
	if err != nil {
		return nil, err
	}
	return NewModelColumns(conn, m.Database, table, m.Prefix)
}
func (m *ModelMapper) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.SlienceMode, "s", false, "Slience mode")
	m.FlagSet().StringVar(&m.Database, "database", "database",
		`database module name. 
	`)
	m.FlagSet().StringVar(&m.QueryID, "id", "common",
		`moder mapper id for actions,queries and viewmodels. 
	`)
	m.FlagSet().StringVar(&m.Location, "location", "modelmappers",
		`default model code location. 
	`)
	m.FlagSet().StringVar(&m.User, "user", "",
		`create form with given user module. 
	`)
	m.FlagSet().StringVar(&m.Prefix, "prefix", "",
		`table field prefix. 
	`)
	crud := m.FlagSet().Bool("crud", false, "Whether create all CRUD codes")
	m.FlagSet().BoolVar(&m.CreateAction, "createaction", false, "Whether create model actions")
	m.FlagSet().BoolVar(&m.CreateForm, "createform", false, "Whether create model forms")
	m.FlagSet().BoolVar(&m.CreateViewModel, "createviewmodel", false, "Whether create model viewmodel class")
	m.FlagSet().BoolVar(&m.WithCreate, "withcreate", false, "Whether create model create code")
	m.FlagSet().BoolVar(&m.WithRead, "withread", false, "Whether create model read code")
	m.FlagSet().BoolVar(&m.WithUpdate, "withupdate", false, "Whether create model update code")
	m.FlagSet().BoolVar(&m.WithDelete, "withdelete", false, "Whether create model delete code")
	m.FlagSet().BoolVar(&m.WithList, "withlist", false, "Whether create model list code")
	m.FlagSet().BoolVar(&m.WithPager, "withpager", false, "Whether create model pager code")
	m.FlagSet().BoolVar(&m.AutoPK, "autopk", false, "auto fill primay key field with unique id")
	m.FlagSet().StringVar(&m.CreatedTimestampField, "createdtimestampfield", "",
		`created timestamp field.model will auto fill this field 
	`)
	m.FlagSet().StringVar(&m.UpdatedTimestampField, "updatedtimestampfield", "",
		`updated timestamp field.model will auto fill this field 
	`)
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
		m.CreateViewModel = true
		m.WithList = true
		m.WithPager = true
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *ModelMapper) Question(a *app.Application, mc *ModelColumns) error {
	var err error
	if m.SlienceMode {
		return nil
	}
	err = QuestionAutoPk.ExecIf(a, !m.AutoPK && mc.CanAutoPK, &m.AutoPK)
	if err != nil {
		return err
	}
	if mc.CreatedTimestampField != nil && m.CreatedTimestampField == "" {
		var result bool
		err = NewQuestionCreatedField(mc.CreatedTimestampField).ExecIf(a, true, &result)
		if err != nil {
			return err
		}
		if result == false {
			mc.CreatedTimestampField = nil
		}
	}
	if mc.UpdatedTimestampField != nil && m.UpdatedTimestampField == "" {
		var result bool
		err = NewQuestionUpdatedField(mc.UpdatedTimestampField).ExecIf(a, true, &result)
		if err != nil {
			return err
		}
		if result == false {
			mc.UpdatedTimestampField = nil
		}
	}
	if mc.IsSinglePrimayKey() && (mc.PrimaryKeys[0].ColumnType == "string" || mc.PrimaryKeys[0].ColumnType == "int") {
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
			m.CreateViewModel = true
		}
		err = QuestionWithCreate.ExecIf(a, !m.WithCreate, &m.WithCreate)
		if err != nil {
			return err
		}
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
	} else {
		err := QuestionWithCreate.ExecIf(a, !m.WithCreate, &m.WithCreate)
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
		err = QuestionCreateAction.ExecIf(a, !m.CreateAction, &m.CreateAction)
		if err != nil {
			return err
		}
	}

	if m.CreateAction && (m.WithCreate || m.WithUpdate || m.WithList) {
		m.CreateForm = true
	}
	err = QuestionCreateForm.ExecIf(a, !m.CreateForm, &m.CreateForm)
	if err != nil {
		return err
	}

	if m.CreateAction {
		m.CreateViewModel = true
	}
	err = QuestionCreateViewModel.ExecIf(a, !m.CreateViewModel, &m.CreateViewModel)
	if err != nil {
		return err
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

	if m.QueryID == "" {
		return ErrUnsuportedIDRequired
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
			return fmt.Errorf("User module code   \"%s\"not found", filepath.Join(mp, m.UserModule.LowerWithParentPath, "init.go"))
		}
	}

	mc, err := m.GetColumn(n.Raw)
	if err != nil {
		return err
	}
	if m.CreatedTimestampField != "" {
		mc.CreatedTimestampField, err = name.New(false, m.CreatedTimestampField)
		if err != nil {
			return err
		}
	}
	if m.UpdatedTimestampField != "" {
		mc.UpdatedTimestampField, err = name.New(false, m.UpdatedTimestampField)
		if err != nil {
			return err
		}
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
	err = m.Render(a, a.Cwd, mp, task, n, qn, mc)
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

func (m *ModelMapper) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name, id *name.Name, mc *ModelColumns) error {
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
	if m.AutoPK && mc.CanAutoPK {
		exists, err = tools.FileExists(filepath.Join(mp, "uniqueid", "uniqueid.go"))
		if err != nil {
			return err
		}
		if !exists {
			err := uniqueid.UniqueIDModule.Exec(a, []string{"-s"})
			if err != nil {
				return err
			}
		}
	}
	filesToRender := map[string]string{}
	exists, err = tools.FileExists(filepath.Join(a.Cwd, mp, n.LowerPath("models"), n.Lower+"queries.go"))
	if err != nil {
		return err
	}
	if !exists {
		filesToRender[filepath.Join(mp, n.LowerPath("models"), n.Lower+"queries.go")] = "modelqueries.go.tmpl"
		filesToRender[filepath.Join(mp, n.LowerPath("models"), n.Lower+"columns.go")] = "modelcolumns.go.tmpl"
		filesToRender[filepath.Join(mp, n.LowerPath("models"), n.Lower+"fields.go")] = "modelfields.go.tmpl"
		filesToRender[filepath.Join(mp, n.LowerPath("models"), n.Lower+".go")] = "model.go.tmpl"

	} else {
		a.Printf("File \"%s\" exists.\nSkip model code installation.", filepath.Join(a.Cwd, mp, n.LowerPath("models"), n.Lower+"queries.go"))
	}
	if m.WithList || m.WithCreate || (mc.HasPrimayKey() && (mc.PrimaryKeys[0].ColumnType == "string" || mc.PrimaryKeys[0].ColumnType == "int") && (m.WithDelete || m.WithUpdate || m.WithRead)) {
		if m.CreateForm {
			filesToRender[filepath.Join(mp, n.LowerPath(id.Lower, "forms"), n.Lower+"form.go")] = "modelform.go.tmpl"
		}
		if m.CreateAction {
			filesToRender[filepath.Join(mp, n.LowerPath(id.Lower, "actions"), n.Lower+"action.go")] = "modelaction.go.tmpl"
		}
		if m.CreateViewModel {
			filesToRender[filepath.Join(mp, n.LowerPath(id.Lower, "viewmodels"), n.Lower+"viewmodel.go")] = "modelviewmodel.go.tmpl"
		}
	}

	data := map[string]interface{}{
		"Name":      n,
		"ID":        id,
		"Columns":   mc,
		"Module":    modelmodule,
		"Confirmed": m,
	}
	return task.RenderFiles(filesToRender, data)
}

var ModelMapperModule = &ModelMapper{}
