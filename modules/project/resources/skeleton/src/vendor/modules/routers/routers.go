package routers

import (
	"modules/app"
	"net/http"

	//"modules/actions"
	"github.com/herb-go/herb/file/simplehttpserver"
	"github.com/herb-go/herb/middleware/router"
	"github.com/herb-go/herb/middleware/router/httprouter"
	"github.com/herb-go/herb/render"
	"github.com/herb-go/util"
)

func ActionSuccess(w http.ResponseWriter, r *http.Request) {
	render.MustJSON(w, "success", 200)
}

func New() router.Router {
	var Router = httprouter.New()

	//Only host assests folder if folder exisits
	if app.Assets.URLPrefix != "" {
		Router.StripPrefix(app.Assets.URLPrefix).
			Use(AssestsMiddlewares()...).
			HandleFunc(simplehttpserver.ServeFolder(util.Resources(app.Assets.Location)))
	}
	var RouterAPI = newAPIRouter()
	Router.StripPrefix("/api").
		Use(APIMiddlewares()...).
		Handle(RouterAPI)

	//var RouterHTML = newHTMLRouter()
	//Router.StripPrefix("/page").Use(HTMLMiddlewares()...).Handle(RouterHTML)

	//Router.GET("/").Use(HTMLMiddlewares()...).HandleFunc(actions.IndexAction)

	return Router
}
