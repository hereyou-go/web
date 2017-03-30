package web

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"github.com/hereyou-go/web/http"
)

// Handler 定义的路由处理逻辑，它将生成一个 HandlerFunc 对象。
type Handler interface{}
// HandlerFunc 是经过预处理的 Handler 包裹，并作为请求的结果返回或进一步处理。
type HandlerFunc func(Context) (int, View)
// ServiceFunc 用于获取一个参数（可能改名）
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
func routerIsReturnInt(tp reflect.Type) bool{
	switch tp.Kind() {
		case reflect.Int:
		fallthrough 
		case reflect.Int16:
		fallthrough
		case reflect.Int32:
		fallthrough
		case reflect.Int64:
		fallthrough
		case reflect.Uint:
		fallthrough
		case reflect.Uint16:
		fallthrough
		case reflect.Uint32:
		fallthrough
		case reflect.Uint64:
		return true
	}
	return false
}
//
//var routerViewType = reflect.TypeOf((*View)(nil)).Elem()
//
//func routerIsReturnView(tp reflect.Type) bool{
//	return tp.Implements(routerViewType)
//}
func routerToView(ctx Context, value interface{}) View {
	view, ok := value.(View)
	if ok {

	} else if s, ok := value.(string); ok {
		if strings.HasPrefix(s, "view:") {
			view = ctx.View(s[5:])
		} else {
			view = ctx.Content(s)
		}
	} else {
		panic(fmt.Errorf("unsupport returns value: %+v ", value))
	}
	return view
}

func routerToInt(result interface{}) int  {
	switch val:= result.(type) {
		case int:
		return val 
		case int16:
		return int(val)
		case int32:
		return int(val)
		case int64:
		return int(val)
		case uint:
		return int(val)
		case uint16:
		return int(val)
		case uint32:
		return int(val)
		case uint64:
		return int(val)
	}
	return -1
}

func routerBuildHandle(method reflect.Value, params []ServiceFunc) HandlerFunc {

	exec := func(ctx Context) []reflect.Value{
		arr := make([]reflect.Value, len(params))
		for i := 0; i < len(params); i++ {
			arr[i] = params[i](ctx)
		}
		return method.Call(arr)
	}

	handlerType := method.Type()
	numOut := handlerType.NumOut()
	if numOut > 2 {
		panic(fmt.Sprintf("Handler返回值数量不能大于2个:%v", numOut))
	} else if numOut == 2 {
		//number 
		if !routerIsReturnInt(handlerType.Out(0)) {
			panic(fmt.Sprintf("Handler返回值的第1个参数必须是个整数类型:%v", handlerType.Out(0)))
		}
		
		// p2:=handlerType.Out(2)
		// switch p2.Kind() {
		// 	case reflect.Struct:
		// 	case reflect.String:
		// 	case reflect.Interface:
		// }
		return func(ctx Context) (int, View) {
			returns := exec(ctx)
			return routerToInt(returns[0].Interface()), routerToView(ctx, returns[1].Interface())
		}
		
	} else if numOut == 1 {
		if routerIsReturnInt(handlerType.Out(0)) {
			return func(ctx Context) (int, View) {
				returns := exec(ctx)
				return routerToInt(returns[0].Interface()), ctx.Empty()
			}
		}
		return func(ctx Context) (int, View) {
			returns := exec(ctx)
			return 200, routerToView(ctx, returns[0].Interface())
		}
	}

	return func(ctx Context) (int, View) {
		exec(ctx)
		return 200, ctx.Empty()
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
