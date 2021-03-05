package scanworkerss

import (
	"fmt"
	"os"
	"path/filepath"

	workertools "github.com/herb-go/worker/tools"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

var QuestionConfirmBackup = tools.NewTrueOrFalseQuestion("Have you backed up or committed your code?")

type ScanWorkers struct {
	app.BasicModule
}

func (m *ScanWorkers) ID() string {
	return "github.com/herb-go/herb-go/modules/scanworkers"
}

func (m *ScanWorkers) Cmd() string {
	return "scanworkers"
}

func (m *ScanWorkers) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s scanworkers.
	scan workers and create workder files.
	File below will be created:
	src/vendor/modules/[module name]/workers.autogenerated.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *ScanWorkers) Desc(a *app.Application) string {
	return "Scan workers and create workder files"
}
func (m *ScanWorkers) Group(a *app.Application) string {
	return "Worker"
}
func (m *ScanWorkers) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *ScanWorkers) Question(a *app.Application) error {
	return nil
}
func (m *ScanWorkers) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}

	if len(args) != 0 {
		a.PrintModuleHelp(m)
		return nil
	}

	mp, err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}
	var confirmed bool
	err = QuestionConfirmBackup.ExecIf(a, true, &confirmed)
	if err != nil {
		return err
	}
	if confirmed == false {
		return nil
	}
	err = QuestionConfirmBackup.ExecIf(a, true, &confirmed)
	if err != nil {
		return err
	}
	if confirmed == false {
		return nil
	}
	if err != nil {
		return err
	}
	c := workertools.NewContext()
	c.GomodFolder = filepath.Join(a.Cwd, "src")
	c.Root = filepath.Join(a.Cwd, mp)
	c.Writer = os.Stdout
	c.MustLoadOverseers("modules/overseers")
	c.MustCheckFolder(filepath.Join(a.Cwd, "src"))
	c.MustRenderAndWrite()
	return nil
}

var ScanWorkersModule = &ScanWorkers{}

func init() {
	app.Register(ScanWorkersModule)
}
