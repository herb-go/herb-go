package app

import (
	"sync/atomic"

	forwarded "github.com/herb-go/herb/middleware/forwarded"
	"github.com/herb-go/herb/middleware/misc"
	"github.com/herb-go/herb/service/httpservice"
	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
)

//HTTPConfig app http config struct
type HTTPConfig struct {
	Forwarded forwarded.Middleware
	Config    httpservice.Config
	Headers   misc.Headers
}

//HTTP app http config
var HTTP = &HTTPConfig{
	Forwarded: forwarded.Middleware{},
	Config:    httpservice.Config{},
	Headers:   misc.Headers{},
}

var syncHTTP atomic.Value

//StoreHTTP atomically store http config
func (a *appSync) StoreHTTP(c *HTTPConfig) {
	syncHTTP.Store(c)
}

//LoadHTTP atomically load http config
func (a *appSync) LoadHTTP() *HTTPConfig {
	v := syncHTTP.Load()
	if v == nil {
		return nil
	}
	return v.(*HTTPConfig)
}

func init() {
	config.RegisterLoaderAndWatch(util.ConfigFile("/http.toml"), func(configpath source.Source) {
		util.Must(tomlconfig.Load(configpath, HTTP))
		Sync.StoreHTTP(HTTP)
		util.SetWarning("Forwarded", HTTP.Forwarded.Warnings()...)
	})
}
