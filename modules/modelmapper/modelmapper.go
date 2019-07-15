package modelmapper

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

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
func (m *ModelMapper) Init(a *app.Application, args *[]string) error {
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
func (m *ModelMapper) Question(a *app.Application) error {
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
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}

	task := tools.NewTask(filepath.Join(app, "/modules/config/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Config  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *ModelMapper) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {
	var configgopath string

	filesToRender := map[string]string{
		filepath.Join("system", "constants", n.LowerWithParentDotSeparated+".toml"): "config.toml.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"):               configgopath,
	}
	return task.RenderFiles(filesToRender, n)
}

var ModelMapperModule = &ModelMapper{}

func init() {
	app.Register(ModelMapperModule)
}
