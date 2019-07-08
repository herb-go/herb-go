package form

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)
var WithActionQuestion = tools.NewQuestion().
	SetDescription("Do you want to create form validating action").
	AddAnswer("y", "Yes", true).
	AddAnswer("n", "No", false)


type Form struct {
	app.BasicModule
	AutoConfirm bool
	WithAction bool
}

func (m *Form) ID() string {
	return "github.com/herb-go/herb-go/modules/form"
}

func (m *Form) Cmd() string {
	return "form"
}

func (m *Form) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s form  <options> <name>.
	Create form file.
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Form) Desc(a *app.Application) string {
	return "Create form"
}
func (m *Form) Group(a *app.Application) string {
	return "Web"
}
func (m *Form) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.AutoConfirm, "y", false, "Whether auto confirm")
	m.FlagSet().BoolVar(&m.WithAction, "withaction", false, "Whether create form action")

	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Form) Question(a *app.Application) error {
	return WithActionQuestion.ExecIf(a,!m.WithAction,&m.WithAction)
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
	mp,err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}

	task := tools.NewTask(filepath.Join(app, "/modules/form/resources"), a.Cwd)
	err = m.Question(a)
	if err != nil {
		return err
	}
	err = m.Render(a, a.Cwd,mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Form  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *Form) Render(a *app.Application, appPath string,mp string, task *tools.Task, n *name.Name) error {
	filesToRender := map[string]string{
		filepath.Join(mp, n.LowerPath("forms"), n.Lower+"form.go"):           "form.go.tmpl",
	}
	if m.WithAction{
		filesToRender[filepath.Join(mp, n.LowerPath("actions"), n.Lower+"action.go")]="action.go.tmpl"
	}
	return task.RenderFiles(filesToRender, n)
}

var FormModule = &Form{}

func init() {
	app.Register(FormModule)
}