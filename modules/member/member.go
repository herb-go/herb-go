package member

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/cache"
	"github.com/herb-go/herb-go/modules/database"
	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/herb-go/modules/session"
	"github.com/herb-go/herb-go/modules/uniqueid"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Member struct {
	app.BasicModule
	InstallSession    bool
	InstallSQLUser    bool
	InstallDatabase   bool
	DatabaseInstalled bool
	UniqueIDInstalled bool
	Cache             string
	CacheName         *name.Name
	SlienceMode       bool
}

type renderData struct {
	Name              *name.Name
	InstallSession    bool
	InstallSQLUser    bool
	CacheName         *name.Name
	DatabaseInstalled bool
}

func (m *Member) ID() string {
	return "github.com/herb-go/herb-go/modules/member"
}

func (m *Member) Cmd() string {
	return "member"
}

func (m *Member) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s member <name>.
Create member module and config files.
Default name is "member".
File below will be created:
	config/<name>.toml
	system/config.examples/<name>.toml
	src/vendor/modules/app/<name>.go
	src/vendor/modules/<name>/<name>.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Member) Desc(a *app.Application) string {
	return "Create member module and config files"
}
func (m *Member) Group(a *app.Application) string {
	return "Auth"
}
func (m *Member) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.SlienceMode, "s", false, "Slience mode")
	m.FlagSet().StringVar(&m.Cache, "cache", "cache", "cache module name")
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Member) Question(a *app.Application, mp string) error {
	err := tools.NewTrueOrFalseQuestion("Do you want to install session module").ExecIf(a, !m.InstallSession && !m.SlienceMode, &m.InstallSession)
	if err != nil {
		return err
	}
	err = tools.NewTrueOrFalseQuestion("Do you want to install sqluser?\nOtherwise you have to install user modules manually.").ExecIf(a, !m.InstallSQLUser && !m.SlienceMode, &m.InstallSQLUser)
	if err != nil {
		return err
	}
	if m.InstallSQLUser {
		result, err := tools.FileExists(filepath.Join(mp, "database", "database.go"))
		if err != nil {
			return err
		}
		if result {
			m.DatabaseInstalled = true
		} else {
			err = tools.NewTrueOrFalseQuestion("Database module not found.\nDo you want to install database module?Otherwise you have to install database modules manually.").ExecIf(a, !m.InstallDatabase && !m.SlienceMode, &m.InstallDatabase)
			if err != nil {
				return err
			}
		}
		result, err = tools.FileExists(filepath.Join(mp, "uniqueid", "uniqueid.go"))
		if err != nil {
			return err
		}
		if result {
			m.UniqueIDInstalled = true
		}
	}
	return nil
}
func (m *Member) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name

	if len(args) == 0 {
		fmt.Println("No member module name given.\"members\" is used")
		n, err = name.New(true, "members")
	} else {
		n, err = name.New(true, args...)
	}
	if err != nil {
		return err
	}

	mp, err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}
	m.CacheName, err = name.New(true, m.Cache)
	if err != nil {
		return err
	}
	result, err := tools.FileExists(filepath.Join(mp, m.CacheName.LowerWithParentPath, "init.go"))
	if err != nil {
		return err
	}
	if !result {
		err = cache.Module.Exec(a, []string{m.Cache})
		if err != nil {
			return err
		}
	}
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}
	err = m.Question(a, mp)
	if err != nil {
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/member/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Member  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *Member) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {
	if m.InstallSession {
		err := session.SessionModule.Exec(a, []string{"-s", filepath.Join(n.Parents, n.Lower, "/session")})
		if err != nil {
			return err
		}
	}
	if m.InstallDatabase {
		err := database.DatabaseModule.Exec(a, []string{"-s"})
		if err != nil {
			return err
		}
		m.DatabaseInstalled = true
	}
	if !m.UniqueIDInstalled {
		err := uniqueid.UniqueIDModule.Exec(a, []string{"-s"})
		if err != nil {
			return err
		}
		m.UniqueIDInstalled = true
	}
	filesToRender := map[string]string{
		filepath.Join(mp, n.LowerPath("init.go")):                             "member.modules.go.tmpl",
		filepath.Join(mp, "middlewares", n.LowerWithParentDotSeparated+".go"): "middleware.go.tmpl",
	}
	data := renderData{
		Name:              n,
		InstallSession:    m.InstallSession,
		InstallSQLUser:    m.InstallSQLUser,
		CacheName:         m.CacheName,
		DatabaseInstalled: m.DatabaseInstalled,
	}
	return task.RenderFiles(filesToRender, data)
}

var MemberModule = &Member{}

func init() {
	app.Register(MemberModule)
}
