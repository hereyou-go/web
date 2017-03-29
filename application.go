package web

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/hereyou-go/web/utils"
)

//Application ...
type Application struct {
	muLock            *sync.Mutex
	isDefault         bool
	middlewares       map[string]Middleware
	globalMiddlewares map[string]Middleware
	paramterMap       map[string]ServiceFunc
	handlers          []HTTPHandler
	attrs             utils.Attribute
	viewEngine        ViewEngine
	routeTable        *RouteTable
	lang              I18N
}

//Use 设置或注册全局变量、中间件、服务等
func (app *Application) Use(registrations ...interface{}) {
	app.muLock.Lock()
	defer app.muLock.Unlock()

	var err error
	c := len(registrations)
	if c == 1 && registrations[0] != nil {
		obj := registrations[0]
		if lang, ok := obj.(I18N); ok {

			app.lang = lang
		} else if val, ok := obj.(ViewEngine); ok {
			app.viewEngine = val
		} else if val, ok := obj.(HTTPHandler); ok {
			err = val.Init(app)
			if err != nil {
				panic(err)
			}
			app.handlers = append(app.handlers, val)
		} else if val, ok := obj.(*RouterGroup); ok {
			val.buildTo(app.routeTable, app)
		} else {
			err = fmt.Errorf("unsupport :%v", reflect.TypeOf(obj).String())
		}
	} else if c == 2 && registrations[0] != nil {
		if key, ok := registrations[0].(string); ok {
			if value, ok := registrations[1].(string); ok {
				app.attrs.SetProperty(key, value)
			} else {
				app.attrs.Set(key, registrations[1])
			}
			return
		} else {
			err = fmt.Errorf("invalid parameters :%v", registrations)
		}
	} else {
		err = fmt.Errorf("invalid parameters :%v", registrations)
	}

	if err == nil {
		return
	}

	panic(err)
}

func (app *Application) Items() utils.Attribute {
	return app.attrs
}

func (app *Application) Attr(key string, defaultValue ...string) string {
	if val := app.attrs.GetString(key); val != "" {
		return val
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}
