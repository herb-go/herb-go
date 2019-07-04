package views

import (
	"github.com/herb-go/herb/render"

	"github.com/herb-go/util"
)

//Modulename module name used in initing and debuging.
const Modulename = "200View"

//Render html templete render
var Render = render.New()

//ViewsInitiator views initiator
var ViewsInitiator func()

func init() {
	util.RegisterModule(Modulename, func() {
		ViewsInitiator()
	})
}
