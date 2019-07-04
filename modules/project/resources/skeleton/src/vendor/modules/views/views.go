package views

import (
	"modules/app"

	"github.com/herb-go/herb/render"
)

//SkinPath website skin path
func SkinPath() string {
	return app.Assets.URLPrefix + "/skin/"
}

//NewRenderData create new render data
func NewRenderData(title string, additionalRenderData ...render.Data) render.Data {
	data := render.Data{}
	data.Set("Name", app.Website.Name)
	data.Set("SkinPath", SkinPath())
	data.Set("BaseURL", app.HTTP.Config.BaseURL)
	data.Set("MetaKeywords", "")
	data.Set("MetaDescription", app.Website.Description)
	data.Set("Title", title)
	for _, v := range additionalRenderData {
		data.Merge(&v)
	}
	return data
}

//ViewIndex index view
var ViewIndex = Render.GetView("index")
