package project

//AppSkeleton app skeleton map
var AppSkeleton = func(mp string) map[string]string {
	return map[string]string{
		"/.gitignore":                              "/skeleton/.gitignore.example",
		"/appdata/readme.md":                       "/skeleton/appdata/readme.md",
		"/config/readme.md":                        "/skeleton/config/readme.md",
		"/config/development.toml":                 "/skeleton/config/development.toml",
		"/config/loggers.toml":                     "/skeleton/config/loggers.toml",
		"/system/data/readme.md":                   "/skeleton/system/data/readme.md",
		"/test/testconfig/readme.md":               "/skeleton/test/testconfig/readme.md",
		"/src/main.go":                             "/skeleton/src/main.go",
		"/src/errors.go":                           "/skeleton/src/errors.go",
		"/src/build/build-linux.sh":                "/skeleton/src/build/build-linux.sh",
		"/system/constants/time.toml":              "/skeleton/system/constants/time.toml",
		"/system/config.examples/development.toml": "/skeleton/system/config.examples/development.toml",
		"/system/config.examples/loggers.toml":     "/skeleton/system/config.examples/loggers.toml",
		"/resources/readme.md":                     "/skeleton/resources/readme.md",
		mp + "/init.go":                            "/skeleton/src/vendor/modules/init.go",
		mp + "/app/app.sync.go":                    "/skeleton/src/vendor/modules/app/app.sync.go",
		mp + "/app/development.go":                 "/skeleton/src/vendor/modules/app/development.go",
		mp + "/app/loggers.go":                     "/skeleton/src/vendor/modules/app/loggers.go",
		mp + "/app/app_test.go":                    "/skeleton/src/vendor/modules/app/app_test.go",
		mp + "/overseers/init.go":                  "/skeleton/src/vendor/modules/overseers/init.go",
		mp + "/readme.md":                          "/skeleton/src/vendor/modules/readme.md",
		mp + "/go.mod.example":                     "/skeleton/src/vendor/modules/go.mod.example",
		mp + "/app/time.go":                        "/skeleton/src/vendor/modules/app/time.go",
		mp + "/appevents/appevents.go":             "/skeleton/src/vendor/modules/appevents/appevents.go",
		mp + "/appevents/registeredevents.go":      "/skeleton/src/vendor/modules/appevents/registeredevents.go",
		mp + "/loggers/init.go":                    "/skeleton/src/vendor/modules/loggers/init.go",
		mp + "/loggers/init_linux.go":              "/skeleton/src/vendor/modules/loggers/init_linux.go",
		mp + "/test/test.go":                       "/skeleton/src/vendor/modules/test/test.go",
		mp + "/test/init_test.go":                  "/skeleton/src/vendor/modules/test/init_test.go",
		mp + "/test/readme.md":                     "/skeleton/src/vendor/modules/test/readme.md",
		mp + "/test/tests/init.go":                 "/skeleton/src/vendor/modules/test/tests/init.go",
		mp + "/test/tests/init_test.go":            "/skeleton/src/vendor/modules/test/tests/init_test.go",
		mp + "/test/tests/readme.md":               "/skeleton/src/vendor/modules/test/tests/readme.md",
		mp + "/drivers/cache.go":                   "/skeleton/src/vendor/modules/drivers/cache.go",
		mp + "/drivers/marshaler.go":               "/skeleton/src/vendor/modules/drivers/marshaler.go",
		mp + "/drivers/database.go":                "/skeleton/src/vendor/modules/drivers/database.go",
		mp + "/drivers/uniqueid.go":                "/skeleton/src/vendor/modules/drivers/uniqueid.go",
	}
}

//HTTPSkeleton http skeleton map
var HTTPSkeleton = func(mp string) map[string]string {
	return map[string]string{
		"/config/http.toml":                 "/skeleton/config/http.toml",
		"/system/config.examples/http.toml": "/skeleton/config/http.toml",
		"/config/csrf.toml":                 "/skeleton/config/csrf.toml",
		"/system/config.examples/csrf.toml": "/skeleton/config/csrf.toml",
		"/system/constants/assets.toml":     "/skeleton/system/constants/assets.toml",
		"/src/http.go":                      "/skeleton/src/http.go",
		mp + "/app/http.go":                 "/skeleton/src/vendor/modules/app/http.go",
		mp + "/app/csrf.go":                 "/skeleton/src/vendor/modules/app/csrf.go",
		mp + "/app/assets.go":               "/skeleton/src/vendor/modules/app/assets.go",
		mp + "/middlewares/middlewares.go":  "/skeleton/src/vendor/modules/middlewares/middlewares.go",
		mp + "/middlewares/csrf.go":         "/skeleton/src/vendor/modules/middlewares/csrf.go",
		mp + "/routers/api.go":              "/skeleton/src/vendor/modules/routers/api.go",
		mp + "/routers/assests.go":          "/skeleton/src/vendor/modules/routers/assests.go",
		mp + "/routers/routers.go":          "/skeleton/src/vendor/modules/routers/routers.go",
	}
}

// WebsiteSkeleton website skeleton map
var WebsiteSkeleton = func(mp string) map[string]string {
	return map[string]string{
		"/system/constants/website.toml":  "/skeleton/system/constants/website.toml",
		"/resources/errorpages/404.html":  "/skeleton/resources/errorpages/404.html",
		"/resources/errorpages/500.html":  "/skeleton/resources/errorpages/500.html",
		mp + "/actions/website.go":        "/skeleton/src/vendor/modules/actions/website.go",
		mp + "/middlewares/errorpages.go": "/skeleton/src/vendor/modules/middlewares/errorpages.go",
		mp + "/app/website.go":            "/skeleton/src/vendor/modules/app/website.go",
		mp + "/routers/html.go":           "/skeleton/src/vendor/modules/routers/html.go",
		mp + "/views/init.go":             "/skeleton/src/vendor/modules/views/init.go",
		mp + "/views/views.go":            "/skeleton/src/vendor/modules/views/views.go",
	}
}

// JetEngineSkeleton jet engine skeleton map
var JetEngineSkeleton = func(mp string) map[string]string {
	return map[string]string{
		"/resources/template.jet/views.toml":       "/skeleton/resources/template.jet/views.toml",
		"/resources/template.jet/layouts/main.jet": "/skeleton/resources/template.jet/layouts/main.jet",
		"/resources/template.jet/views/index.jet":  "/skeleton/resources/template.jet/views/index.jet",
		mp + "/views/jetengine.go":                 "/skeleton/src/vendor/modules/views/jetengine.go",
	}
}

// TmplEngineSkeleton tmpl engine skeleton map
var TmplEngineSkeleton = func(mp string) map[string]string {
	return map[string]string{
		"/resources/template.tmpl/views.toml":        "/skeleton/resources/template.tmpl/views.toml",
		"/resources/template.tmpl/layouts/main.tmpl": "/skeleton/resources/template.tmpl/layouts/main.tmpl",
		"/resources/template.tmpl/views/index.tmpl":  "/skeleton/resources/template.tmpl/views/index.tmpl",
		mp + "/views/tmplengine.go":                  "/skeleton/src/vendor/modules/views/tmplengine.go",
	}
}
