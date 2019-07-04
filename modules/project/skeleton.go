package project

var AppSkeleton = map[string]string{
	"/.gitignore":                                       "/skeleton/.gitignore",
	"/appdata/readme.md":                                "/skeleton/appdata/readme.md",
	"/config/readme.md":                                 "/skeleton/config/readme.md",
	"/config/development.toml":                          "/skeleton/config/development.toml",
	"/resources/data/readme.md":                         "/skeleton/resources/data/readme.md",
	"/src/test/config/readme.md":                        "/skeleton/src/test/config/readme.md",
	"/src/main.go":                                      "/skeleton/src/main.go",
	"/src/drivers.go":                                   "/skeleton/src/drivers.go",
	"/src/errors.go":                                    "/skeleton/src/errors.go",
	"/src/vendor/modules/app/development.go":            "/skeleton/src/vendor/modules/app/development.go",
	"/src/vendor/modules/app/loadconfigs.go":            "/skeleton/src/vendor/modules/app/loadconfigs.go",
	"/src/vendor/modules/app/time.go":                   "/skeleton/src/vendor/modules/app/time.go",
	"/src/vendor/modules/appevents/appevents.go":        "/skeleton/src/vendor/modules/appevents/appevents.go",
	"/src/vendor/modules/appevents/registeredevents.go": "/skeleton/src/vendor/modules/appevents/registeredevents.go",
	"/system/constants/time.toml":                       "/skeleton/system/constants/time.toml",
	"/system/config.examples/development.toml":          "/skeleton/system/config.examples/development.toml",
}

var HTTPSkeleton = map[string]string{
	"/config/http.toml":                              "/skeleton/config/http.toml",
	"/system/config.examples/http.toml":              "/skeleton/config/http.toml",
	"/config/csrf.toml":                              "/skeleton/config/csrf.toml",
	"/system/config.examples/csrf.toml":              "/skeleton/config/csrf.toml",
	"/system/constants/assets.toml":                  "/skeleton/system/constants/assets.toml",
	"/src/http.go":                                   "/skeleton/src/http.go",
	"/src/vendor/modules/app/http.go":                "/skeleton/src/vendor/modules/app/http.go",
	"/src/vendor/modules/app/csrf.go":                "/skeleton/src/vendor/modules/app/csrf.go",
	"/src/vendor/modules/app/assets.go":              "/skeleton/src/vendor/modules/app/assets.go",
	"/src/vendor/modules/messages/forms.go":          "/skeleton/src/vendor/modules/messages/forms.go",
	"/src/vendor/modules/messages/messages.go":       "/skeleton/src/vendor/modules/messages/messages.go",
	"/src/vendor/modules/middlewares/middlewares.go": "/skeleton/src/vendor/modules/middlewares/middlewares.go",
	"/src/vendor/modules/middlewares/csrf.go":        "/skeleton/src/vendor/modules/middlewares/csrf.go",
	"/src/vendor/modules/routers/api.go":             "/skeleton/src/vendor/modules/routers/api.go",
	"/src/vendor/modules/routers/assests.go":         "/skeleton/src/vendor/modules/routers/assests.go",
	"/src/vendor/modules/routers/routers.go":         "/skeleton/src/vendor/modules/routers/routers.go",
}

var WebsiteSkeleton = map[string]string{
	"/system/constants/website.toml":                "/skeleton/system/constants/website.toml",
	"/resources/errorpages/404.html":                "/skeleton/resources/errorpages/404.html",
	"/resources/errorpages/500.html":                "/skeleton/resources/errorpages/500.html",
	"/src/vendor/modules/actions/website.go":        "/skeleton/src/vendor/modules/actions/website.go",
	"/src/vendor/modules/middlewares/errorpages.go": "/skeleton/src/vendor/modules/middlewares/errorpages.go",
	"/src/vendor/modules/app/website.go":            "/skeleton/src/vendor/modules/app/website.go",
	"/src/vendor/modules/routers/html.go":           "/skeleton/src/vendor/modules/routers/html.go",
	"/src/vendor/modules/views/init.go":             "/skeleton/src/vendor/modules/views/init.go",
	"/src/vendor/modules/views/views.go":            "/skeleton/src/vendor/modules/views/views.go",
}

var JetEngineSkeleton = map[string]string{
	"/resources/template.jet/views.toml":       "/skeleton/resources/template.jet/views.toml",
	"/resources/template.jet/layouts/main.jet": "/skeleton/resources/template.jet/layouts/main.jet",
	"/resources/template.jet/views/index.jet":  "/skeleton/resources/template.jet/views/index.jet",
	"/src/vendor/modules/views/jetengine.go":   "/skeleton/src/vendor/modules/views/jetengine.go",
}

var TmplEngineSkeleton = map[string]string{
	"/resources/template.tmpl/views.toml":        "/skeleton/resources/template.tmpl/views.toml",
	"/resources/template.tmpl/layouts/main.tmpl": "/skeleton/resources/template.tmpl/layouts/main.tmpl",
	"/resources/template.tmpl/views/index.tmpl":  "/skeleton/resources/template.tmpl/views/index.tmpl",
	"/src/vendor/modules/views/tmplengine.go":    "/skeleton/src/vendor/modules/views/tmplengine.go",
}
