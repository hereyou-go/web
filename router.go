package web

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/hereyou-go/web/http"
)

type Handler interface{}
type HandlerFunc func(Context) interface{}
type ServiceFunc func(Context) reflect.Value

type Router struct {
	method      http.HttpMethod
	pattern     string
	middlewares []string
	handler     Handler
	handle      HandlerFunc
	controller  interface{}
	// group       *RouterGroup
	// isGroup     bool
	// parent      *Router
	// items       []*Router
}

func routerGetParams(handlerType reflect.Type, app *Application, instance interface{}) []ServiceFunc {
	params := make([]ServiceFunc, handlerType.NumIn())
	for i := 0; i < handlerType.NumIn(); i++ {
		//ptyp := handlerType.In(i)
		typeName := handlerType.In(i).String()
		fn, ok := app.paramterMap[typeName]
		if !ok {
			switch typeName {
			case "*web.Application":
				param := reflect.ValueOf(app)
				fn = func(ctx Context) reflect.Value {
					return param
				}
			case "web.Context":
				fn = func(ctx Context) reflect.Value {
					return reflect.ValueOf(ctx)
				}
			case "*web.Request":
				fn = func(ctx Context) reflect.Value {
					return reflect.ValueOf(ctx.Request())
				}
			case "*web.Response":
				fn = func(ctx Context) reflect.Value {
					return reflect.ValueOf(ctx.Response())
				}
			case "*http.Request":
				fn = func(ctx Context) reflect.Value {
					return reflect.ValueOf(ctx.Request().RawRequest())
				}
			case "http.ResponseWriter":
				fn = func(ctx Context) reflect.Value {
					return reflect.ValueOf(ctx.Response().Writer())
				}
			default:
				for _, val := range app.attrs {
					vt := reflect.TypeOf(val)
					if vt.String() == typeName {
						param := reflect.ValueOf(val)
						fn = func(ctx Context) reflect.Value {
							return param
						}
					}
				}
				break
			}
			// 如果给定对象为实例方法，则传入实例
			// TODO: 观察是否会有参数冲突
			if fn == nil && i == 0 && instance != nil {
				typ := reflect.TypeOf(instance)
				if handlerType.In(i) == typ {
					param := reflect.ValueOf(instance)
					fn = func(ctx Context) reflect.Value {
						return param
					}
				}
			}
			if fn == nil {
				panic(fmt.Sprintf("参数不支持或未定义:%v", typeName))
			}
			app.paramterMap[typeName] = fn

		}
		params[i] = fn
	}
	return params
}

func routerBuildHandle(method reflect.Value, params []ServiceFunc) HandlerFunc {
	return func(ctx Context) interface{} {
		arr := make([]reflect.Value, len(params))
		for i := 0; i < len(params); i++ {
			arr[i] = params[i](ctx)
		}
		returns := method.Call(arr)
		if len(returns) > 0 {
			return returns[0].Interface()
		}
		return nil
	}
}

func buildHandler(app *Application, handler Handler, instance interface{}) (handle HandlerFunc) {
	if method, ok := handler.(reflect.Method); ok {
		handle = routerBuildHandle(method.Func, routerGetParams(method.Func.Type(), app, instance))
	} else {
		handlerType := reflect.TypeOf(handler)
		switch handlerType.Kind() {
		case reflect.Func:
			handle = routerBuildHandle(reflect.ValueOf(handler), routerGetParams(handlerType, app, instance))
		default:
			fmt.Println(handlerType.Kind())
		}
	}
	return
}

var patternRegexp = regexp.MustCompile(`\{\s*(\w+)\s*(\??)\s*\}`)

func compilePattern(pattern string) (string, []string) {
	var keys []string
	var index []int
	var matched []string
	for {
		index = patternRegexp.FindStringIndex(pattern)
		if len(index) != 2 {
			break
		}
		matched = patternRegexp.FindStringSubmatch(pattern)

		if len(matched) != 3 {
			break
		}
		keys = append(keys, matched[1])
		partten := "(?P<" + matched[1] + ">\\w+)"
		if matched[2] == "?" {
			partten += "?"
			//TODO:如果后面有空格的话将会判断不到
			//t.Fatalf("****：%v----%v", index[1]+1, len(pattern))
			if len(pattern) == index[1]+1 && pattern[index[1]] == '/' {
				index[1] = index[1] + 1
				partten += "/?"
			}
		}
		pattern = pattern[:index[0]] + partten + pattern[index[1]:]
		// t.Fatalf("****：%v", len(matched))
		//break
	}
	return pattern, keys
}
