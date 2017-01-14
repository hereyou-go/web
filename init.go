package web

import (
	"container/list"
	"sync"

	"github.com/one-go/web/utils"
)

//The Application interface is ...
// type Application interface {
// }

var defaultApp *Application

func init() {
	defaultApp = New()

	//defaultApp.Use(RouteHandler)
}

//New returns a new Application.
func New() *Application {
	return &Application{
		muLock:      new(sync.Mutex),
		isDefault:   true,
		attrs:       make(utils.Attribute),
		handlers:    make([]HTTPHandler, 0),
		paramterMap: make(map[string]ServiceFunc),
		routeTable: &RouteTable{
			routes: list.New(),
		},
	}
}

//Default returns a default(global) Engine instance.
func Default() *Application {
	return defaultApp
}

//AddMiddleware method add a Middleware to Application.
func (app *Application) AddMiddleware(name string, ware Middleware, global ...bool) {
	if global != nil && len(global) >= 1 && global[0] == true {
		app.globalMiddlewares[name] = ware
	} else {
		app.middlewares[name] = ware
	}

}
