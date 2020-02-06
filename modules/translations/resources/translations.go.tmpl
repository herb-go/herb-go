package app

import (
	"sync/atomic"

	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
	"github.com/herb-go/util/config/translateconfig"
)

//Translations app translations config
var Translations = &translateconfig.Config{}

var syncTranslations atomic.Value

//StoreTranslations atomically store translations config
func (a *appSync) StoreTranslations(c *translateconfig.Config) {
	syncTranslations.Store(c)
}

//LoadTranslations atomically load translations config
func (a *appSync) LoadTranslations() *translateconfig.Config {
	v := syncTranslations.Load()
	if v == nil {
		return nil
	}
	return v.(*translateconfig.Config)
}


func init() {
	config.RegisterLoader(util.SystemFile("messages", "translations.toml"), func(configpath source.Source) {
		util.Must(tomlconfig.Load(configpath, Translations))
		Sync.StoreTranslations(Translations)
	})
}
