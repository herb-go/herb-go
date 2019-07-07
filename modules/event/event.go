package router

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Event struct {
	app.BasicModule
	AutoConfirm bool
}

func (m *Event) ID() string {
	return "github.com/herb-go/herb-go/modules/event"
}

func (m *Event) Cmd() string {
	return "event"
}

func (m *Event) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s event <name>.
Create event files.
File below will be created:
	src/vendor/modules/appevents/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Event) Desc(a *app.Application) string {
	return "Create app event file"
}
func (m *Event) Group(a *app.Application) string {
	return "Application"
}
func (m *Event) Init(a *app.Application, args *[]string) error {
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
func (m *Event) Question(a *app.Application) error {
	return nil
}
func (m *Event) Exec(a *app.Application, args []string) error {
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

	task := tools.NewTask(filepath.Join(app, "/modules/event/resources"), a.Cwd)

	err = m.Render(a, a.Cwd,mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Event  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *Event) Render(a *app.Application, appPath string,mp string, task *tools.Task, n *name.Name) error {

	filesToRender := map[string]string{
		filepath.Join(mp,"appevents", n.Lower+".go"):           "event.go.tmpl",
	}
	return task.RenderFiles(filesToRender, n)
}

var EventModule = &Event{}

func init() {
	app.Register(EventModule)
}
