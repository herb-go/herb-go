package routers

import (
	"modules/middlewares"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/middleware/router"
	"github.com/herb-go/herb/middleware/router/httprouter"
)

//HTMLMiddlewares middlewares that should used in html requests
var HTMLMiddlewares = func() middleware.Middlewares {
	return middleware.Middlewares{
		middlewares.MiddlewareCsrfSetToken,
	}
}

//RouterHTML html router
var RouterHTML = newHTMLRouter()

func newHTMLRouter() router.Router {
	var Router = httprouter.New()
	//Put your router configure code here
	// Router.GET("/page/:id").
	// 	HandleFunc(actions.PageAction)
	return Router

}
