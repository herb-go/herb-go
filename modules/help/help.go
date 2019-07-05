package help

import (
	"sort"
	"fmt"

	"github.com/herb-go/util/cli/app"
)

type helpGroup struct{
	Group string
	Value string
}

type helpGroups []*helpGroup

func (g *helpGroups)        Len() int{
	return len(*g)
}
func (g *helpGroups)        Less(i, j int) bool{
	return (*g)[i].Group<(*g)[j].Group
}
func (g *helpGroups)        Swap(i, j int){
	t:=(*g)[i]
	(*g)[i]=(*g)[j]
	(*g)[j]=t
}
func (g *helpGroups) Add(group,value string){
	for k:=range *g{
		if (*g)[k].Group==group{
		 	(*g)[k].Value = (*g)[k].Value  +  value
		 	return
		}
	}
	(*g)=append(*g,&helpGroup{Group:group,Value:value})
}


func (g *helpGroups) String() string{
	var result string
	sort.Sort(g)

	for _,v:=range (*g){
		if v.Group!=""{
			result += fmt.Sprintf("[%s]\r\n",v.Group)
		}
		result+=v.Value
	}
	return result
}

func newHelpGroups()*helpGroups{
	return &helpGroups{}
}

type Help struct {
	app.BasicModule
}

func (m *Help) Cmd() string {
	return "help"
}

func (m *Help) ID() string {
	return "github.com/herb-go/herb-go/modules/help"
}


func (m *Help) Help(a *app.Application) string {
	help := fmt.Sprintf("Usage %s help [command]\r\n",a.Config.Cmd)
	help += "Command list:\r\n"
	groups:=newHelpGroups()
	for _, v := range *a.Modules {
		groups.Add(v.Group(a),fmt.Sprintf("  %s - %s\r\n", v.Cmd(), v.Desc(a)))
	}
	help+=groups.String()
	return help
}

func (m *Help) Desc(a *app.Application) string {
	return "Show module help"
}

func (m *Help) Exec(a *app.Application, args []string) error {
	if len(args) != 1 {
		a.PrintModuleHelp(m)
		return nil
	}
	cmd := args[0]
	module := a.Modules.Get(cmd)
	if module == nil {
		a.Printf("Module \"%s\" not found.\n", cmd)
		a.Println(m)
		return nil
	}
	a.PrintModuleHelp(module)
	return nil
}

var Module = &Help{}

func init() {
	app.Register(Module)
}
