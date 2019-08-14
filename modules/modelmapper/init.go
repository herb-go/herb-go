package modelmapper

import "github.com/herb-go/util/cli/app"

func init() {
	app.Register(ModelMapperModule)
	app.Register(UpdateModule)
	app.Register(DataSourceModule)
	app.Register(FormModule)
	app.Register(SelectModule)
}
