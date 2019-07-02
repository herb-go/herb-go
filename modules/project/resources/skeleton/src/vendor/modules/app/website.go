package app

import (
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/commonconfig"
	"github.com/herb-go/util/config/tomlconfig"
)

var Website = commonconfig.WebsiteConfig{}

func init() {
	config.RegisterLoader(util.ConstantsFile("/website.toml"), func(configpath util.FileObject) {
		util.Must(tomlconfig.Load(configpath, &Website))
	})
}
