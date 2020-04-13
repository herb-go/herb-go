package worker

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
	"github.com/herb-go/util/cli/name"
)

type Worker struct {
	app.BasicModule
	WorkerName  string
	Name        *name.Name
	Folder      string
	SlienceMode bool
}

func (m *Worker) ID() string {
	return "github.com/herb-go/herb-go/modules/worker"
}

func (m *Worker) Cmd() string {
	return "worker"
}

func (m *Worker) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s worker.
Create sql worker module and config files.
File below will be created:
	config/worker.toml
	system/confg.examples/worker.toml
	src/vendor/modules/app/worker.go
	src/vendor/modules/worker/init.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Worker) Desc(a *app.Application) string {
	return "Create worker module files."
}
func (m *Worker) Group(a *app.Application) string {
	return "Worker"
}
func (m *Worker) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.SlienceMode, "s", false, "Slience mode")
	m.FlagSet().StringVar(&m.WorkerName, "n", "", "Worker name")
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Worker) Question(a *app.Application) error {
	return nil
}
func (m *Worker) Exec(a *app.Application, args []string) error {
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
	m.Name, err = name.New(false, m.WorkerName)
	if err != nil {
		return err
	}
	if m.WorkerName == "" {
		m.Folder = "worker"
	} else {
		m.Folder = m.Name.Lower
	}
	mp, err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}
	result, err := tools.IsFolder(mp, n.Raw)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("folder " + filepath.Join(mp, n.Raw) + " not found")
	}
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}
	err = InitWorkers(a, a.Cwd, mp, m.SlienceMode)
	if err != nil {
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/worker/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, n, task)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Worker " + n.Pascal + m.Name.Pascal + " created.\n")
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

func (m *Worker) Render(a *app.Application, appPath string, mp string, n *name.Name, task *tools.Task) error {
	var renderData = map[string]interface{}{
		"Name":   n,
		"Module": m,
	}
	filesToRender := map[string]string{
		filepath.Join(mp, n.Raw, "herb-workers", m.Folder, "worker.go"): "herb-worker.go.tmpl",
		filepath.Join(mp, "workers", n.Lower+"."+m.Folder+".go"):        "workers.go.tmpl",
	}
	return task.RenderFiles(filesToRender, renderData)
}

var WorkerModule = &Worker{}

func init() {
	app.Register(WorkerModule)
}
