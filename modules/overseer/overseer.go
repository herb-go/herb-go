package overseers

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
	"github.com/herb-go/util/cli/name"
)

type Overseer struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Overseer) ID() string {
	return "github.com/herb-go/herb-go/modules/overseer"
}

func (m *Overseer) Cmd() string {
	return "overseer"
}

func (m *Overseer) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s overseer.
Create overseer module files.
	File below will be created:
	src/vendor/modules/overseer/init.go
	src/vendor/modules/overseer/[name].go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Overseer) Desc(a *app.Application) string {
	return "Create overseer module files."
}
func (m *Overseer) Group(a *app.Application) string {
	return "Worker"
}
func (m *Overseer) Init(a *app.Application, args *[]string) error {
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
func (m *Overseer) Question(a *app.Application) error {
	return nil
}
func (m *Overseer) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}

	if len(args) != 1 {
		a.PrintModuleHelp(m)
		return nil
	}
	n, err := name.New(false, args...)
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
	err = project.InitOverseers(a, a.Cwd, mp, true)
	if err != nil {
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/overseer/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, n, task)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Overseer " + n.Pascal + " created.\n")
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

func (m *Overseer) Render(a *app.Application, appPath string, mp string, n *name.Name, task *tools.Task) error {
	var renderData = map[string]interface{}{
		"Name":   n,
		"Module": m,
	}
	filesToRender := map[string]string{
		filepath.Join(mp, "overseers", n.Lower+".go"): "overseers.go.tmpl",
	}
	return task.RenderFiles(filesToRender, renderData)
}

var OverseerModule = &Overseer{}

func init() {
	app.Register(OverseerModule)
}
