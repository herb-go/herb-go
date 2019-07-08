package member

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/herb-go/modules/session"
	"github.com/herb-go/herb-go/modules/cache"
	"github.com/herb-go/herb-go/modules/database"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Member struct {
	app.BasicModule
	InstallSession    bool
	InstallSQLUser    bool
	InstallCache      bool
	InstallDatabase bool

	AutoConfirm bool
}

type renderData struct {
	Name *name.Name
	InstallSession    bool
	InstallSQLUser    bool
	InstallCache      bool
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
	m.FlagSet().BoolVar(&m.AutoConfirm, "y", false, "Whether auto confirm")
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Member) Question(a *app.Application) error {
	err:=tools.NewTrueOrFalseQuestion("Do you want to install session module").ExecIf(a,!m.InstallSession,&m.InstallSession)
	if err!=nil{
		return err
	}
	err=tools.NewTrueOrFalseQuestion("Do you want to add member cache code?Otherwise you have to install member cache manually").ExecIf(a,!m.InstallCache,&m.InstallCache)
	if err!=nil{
		return err
	}
		err=tools.NewTrueOrFalseQuestion("Do you want to install sqluser?Otherwise you have to install user modules manually.").ExecIf(a,!m.InstallSQLUser,&m.InstallSQLUser)
	if err!=nil{
		return err
	}
	if m.InstallSQLUser{
	err=tools.NewTrueOrFalseQuestion("Database module not found.\nDo you want to install database module?Otherwise you have to install user modules manually.").ExecIf(a,!m.InstallDatabase,&m.InstallDatabase)
	if err!=nil{
		return err
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
		fmt.Println("No member module name given.\"member\" is used")
		n ,err= name.New(true,"member")
	}else{
		n, err = name.New(true, args...)
	}
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
	err=m.Question(a)
	if err!=nil{
		return err
	}
	task := tools.NewTask(filepath.Join(app, "/modules/member/resources"), a.Cwd)

	err = m.Render(a, a.Cwd,mp, task, n)
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
	ok, err := task.ConfirmIf(a, !m.AutoConfirm)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return task.Exec()

}

func (m *Member) Render(a *app.Application, appPath string,mp string, task *tools.Task, n *name.Name) error {
	if m.InstallSession{
		session.SessionModule.AutoConfirm=true
		err:=session.SessionModule.Exec(a,[]string{filepath.Join(n.Parents,n.Lower ,"/session")})
		if err!=nil{
			return err
		}
	}
    if m.InstallCache{
		cache.Module.AutoConfirm=true
		err:=cache.Module.Exec(a,[]string{filepath.Join(n.Parents,n.Lower ,"/cache")})
		if err!=nil{
			return err
		}
	}
	if m.InstallDatabase{
		database.DatabaseModule.AutoConfirm=true
		err:=database.DatabaseModule.Exec(a,[]string{})
		if err!=nil{
			return err
		}
	}
	filesToRender := map[string]string{
		filepath.Join(mp,n.LowerPath(n.Lower+".go")):"member.modules.go.tmpl",
		filepath.Join(mp,"middlewares", n.LowerWithParentDotSeparated+".go"):"middleware.go.tmpl",
	}
	data:=renderData{
		Name:n,
	InstallSession :m.InstallSession,
	InstallSQLUser :m.InstallSQLUser,
	InstallCache:m.InstallCache,
	DatabaseInstalled:m.InstallDatabase,
	}
	return task.RenderFiles(filesToRender, data)
}

var MemberModule = &Member{}

func init() {
	app.Register(MemberModule)
}