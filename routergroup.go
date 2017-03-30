package web

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/hereyou-go/web/http"
	"github.com/hereyou-go/logs"
)

type RouterGroup struct {
	pattern     string
	middlewares []string
	parent      *RouterGroup
	routers     []*Router
	groups      []*RouterGroup
}

func NewGroup(pattern string, middlewares ...string) *RouterGroup {
	group := &RouterGroup{
		pattern:     pattern,
		middlewares: middlewares,
		routers:     make([]*Router, 0),
		groups:      make([]*RouterGroup, 0),
	}
	return group
}

func routerGroupMergeMiddlewares(dst []Middleware, app *Application, middlewares []string) ([]Middleware, error) {
	for _, name := range middlewares {
		if ware,ok:=app.Attr(name).(Middleware);ok{
			dst = append(dst, ware)
		}else {
			return dst, logs.NewError("","指定 Middleware 未定义，使用 app.use 进行注册。name:%v => %v", name,app.Attr(name))
		}
	}
	return  dst,nil
}

func routerGroupMergePattern(group *RouterGroup, app *Application) (string,[]Middleware, error) {
	suffix := strings.Trim(group.pattern, " ")
	if group.parent != nil {
		prefix,wares,err := routerGroupMergePattern(group.parent, app)
		if err!=nil{
			return "",nil,err
		}

		if prefix == "/" {
			prefix = "" //如果上级是默认规则，则去掉
		}
		wares,err=routerGroupMergeMiddlewares(wares,app,group.middlewares)
		if err!=nil{
			return "",nil,err
		}
		return prefix + suffix, wares, nil
	}
	wares:=make([]Middleware,0)
	wares,err:=routerGroupMergeMiddlewares(wares,app,group.middlewares)
	return suffix,wares,err
}

func (group *RouterGroup) buildTo(table *RouteTable, app *Application) error {
	for _, g := range group.groups {
		g.parent = group
		g.buildTo(table, app)
	}
	groupPattern,wares,err := routerGroupMergePattern(group,app)
	if err!=nil{
		return err
	}
	if groupPattern == "/" {
		groupPattern = "" //去掉默认规则
	}
	for _, r := range group.routers {
		pattern, keys := compilePattern(groupPattern + strings.Trim(r.pattern, " "))
		handler := buildHandler(app, r.handler, r.controller)
		rwares:=make([]Middleware,len(wares))
		copy(rwares,wares)
		rwares,err:=routerGroupMergeMiddlewares(rwares,app,r.middlewares)
		if err!=nil{
			return err
		}
		table.Register(r.method, regexp.MustCompile(pattern), keys, handler, 0, false,rwares)
	}
	return nil
}

func (group *RouterGroup) Route(method http.HttpMethod, pattern string, handler Handler, middlewares []string) *Router {
	router := &Router{
		method:      method,
		pattern:     pattern,
		handler:     handler,
		middlewares: middlewares,
	}
	group.routers = append(group.routers, router)
	return router
}

func routerGroupResolveRoute(group *RouterGroup, controller interface{}) {
	typ := reflect.TypeOf(controller)
	styp := typ
	if typ.Kind() == reflect.Ptr {
		styp = typ.Elem()
	}
	var methods map[string]reflect.Method
	methods = make(map[string]reflect.Method)

	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		methods[strings.ToLower(method.Name)] = method
	}
	for i := 0; i < styp.NumField(); i++ {
		f := styp.Field(i)
		method, ok := methods[strings.ToLower(f.Name)]
		if !ok {
			continue
		}

		pattern := f.Tag.Get("route")
		if index := strings.Index(pattern, ":/"); index > -1 {
			httpMethods := strings.Split(pattern[:index], "|")
			pattern = pattern[index+1:]
			fmt.Printf("%d %v = %v , %v   %v\n", i, f.Name, pattern, method.Name, httpMethods) //, f.Type(), f.Interface()
			var httpMethod http.HttpMethod
			isset := false
			for _, s := range httpMethods {
				m := http.ParseHttpMethod(s)
				if isset {
					httpMethod |= m
				} else {
					httpMethod = m
					isset = true
				}
			}
			var wares []string
			group.Route(httpMethod, pattern, method,wares).controller = controller
		}
	}
}

func (group *RouterGroup) AppendController(controller interface{}) *RouterGroup {
	routerGroupResolveRoute(group, controller)
	return group
}

//============== APIs ==============

func (group *RouterGroup) Get(pattern string, handler Handler, middlewares ...string) *Router {
	return group.Route(http.GET, pattern, handler, middlewares)
}

func (group *RouterGroup) Group(pattern string, controller interface{}, middlewares ...string) *RouterGroup {
	sub := NewGroup(pattern, middlewares...).AppendController(controller)
	//sub.middlewares = middlewares
	group.groups = append(group.groups, sub)
	return sub
}
