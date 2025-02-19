package routers

import (
	"elastic_golang/controllers"
	"github.com/beego/beego/v2/server/web" // ✅ Keep only one Beego import
)

func init() {
	web.Router("/", &controllers.MainController{})
	web.Router("/search", &controllers.SearchController{}, "get:SearchHandler")
	web.Router("/autocomplete", &controllers.SearchController{}, "get:AutocompleteHandler")
}
