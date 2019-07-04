package app

import (
	forwarded "github.com/herb-go/herb/middleware/forwarded"
	"github.com/herb-go/herb/middleware/misc"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
	"github.com/herb-go/util/httpserver"
)

//HTTPConfig app http config struct
type HTTPConfig struct {
	Forwarded forwarded.Middleware
	Config    httpserver.Config
	Headers   misc.Headers
}

//HTTP app http config
var HTTP = &HTTPConfig{
	Forwarded: forwarded.Middleware{},
	Config:    httpserver.Config{},
	Headers:   misc.Headers{},
}

func init() {
	config.RegisterLoaderAndWatch(util.ConfigFile("/http.toml"), func(configpath util.FileObject) {
		util.Must(tomlconfig.Load(configpath, HTTP))
		util.SetWarning("Forwarded", HTTP.Forwarded.Warnings()...)
	})
}
