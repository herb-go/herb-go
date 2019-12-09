package app

import (
	"sync/atomic"

	"github.com/herb-go/herbconfig/configuration"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/commonconfig"
	"github.com/herb-go/util/config/tomlconfig"
)

//Website app website settings
var Website = &commonconfig.WebsiteConfig{}

var syncWebsite atomic.Value

//StoreWebsite atomically store website config
func (a *appSync) StoreWebsite(c *commonconfig.WebsiteConfig) {
	syncWebsite.Store(c)
}

//LoadWebsite atomically load website config
func (a *appSync) LoadWebsite() *commonconfig.WebsiteConfig {
	v := syncWebsite.Load()
	if v == nil {
		return nil
	}
	return v.(*commonconfig.WebsiteConfig)
}

func init() {
	config.RegisterLoader(util.ConstantsFile("/website.toml"), func(configpath configuration.Configuration) {
		util.Must(tomlconfig.Load(configpath, Website))
		Sync.StoreWebsite(Website)
	})
}
