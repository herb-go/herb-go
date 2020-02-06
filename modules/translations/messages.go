package translations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Messages struct {
	app.BasicModule
	Modules     string
	modules     []string
	SlienceMode bool
}

func (m *Messages) ID() string {
	return "github.com/herb-go/herb-go/modules/ui.messages"
}

func (m *Messages) Cmd() string {
	return "messages"
}

func (m *Messages) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s constants [language].
Install translated messages.
File below will be created:
	system/messages/[language]/[modules].toml ...
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Messages) Desc(a *app.Application) string {
	return "Install translated messages"
}
func (m *Messages) Group(a *app.Application) string {
	return "Translations"
}
func (m *Messages) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.SlienceMode, "s", false, "Slience mode")
	m.FlagSet().StringVar(&m.Modules, "modules", "", "Moduel name  separated by \",\".Left empty to install all modules")
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Messages) Question(a *app.Application) error {
	return nil
}
func (m *Messages) CheckFolders(a *app.Application, sourcespath string, language string) error {
	result, err := tools.FileExists(sourcespath, language)
	if err != nil {
		return err
	}
	if !result {
		return fmt.Errorf("language \"%s\" not found", language)
	}
	modulelist := []string{}
	modulemap := map[string]bool{}
	filepath.Walk(filepath.Join(sourcespath, language), func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			if strings.HasSuffix(path, ".toml") {
				m := filepath.Base(path[:len(path)-5])
				modulelist = append(modulelist, m)
				modulemap[m] = true
			}
		}
		return nil
	})
	if m.Modules == "" {
		m.modules = modulelist
		return nil
	}
	for _, v := range m.modules {
		if modulemap[v] == false {
			return fmt.Errorf("module message\"%s\" not found", v)
		}
	}
	return nil
}
func (m *Messages) Exec(a *app.Application, args []string) error {
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
	m.modules = strings.Split(m.Modules, ",")
	mp, err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}
	err = m.CheckFolders(a, filepath.Join(app, "/modules/translations/resources/messages"), n.Raw)
	if err != nil {
		return err
	}

	task := tools.NewTask(filepath.Join(app, "/modules/translations/resources"), a.Cwd)
	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Translatition messages  for language \"%s\" created.\n", n.Raw)
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

func (m *Messages) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {
	filesToRender := map[string]string{}
	for _, v := range m.modules {
		filesToRender[filepath.Join("system", "messages", n.Raw, v+".toml")] = "messages/" + n.Raw + "/" + v + ".toml"
	}
	return task.CopyFiles(filesToRender)
}

var MessagesModule = &Messages{}

func init() {
}
