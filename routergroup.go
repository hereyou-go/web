package web

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/one-go/web/http"
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

func routerGroupMergePattern(group *RouterGroup) string {
	suffix := strings.Trim(group.pattern, " ")
	if group.parent != nil {
		prefix := routerGroupMergePattern(group.parent)
		if prefix == "/" {
			prefix = "" //如果上级是默认规则，则去掉
		}
		return prefix + suffix
	}
	return suffix
}

func (group *RouterGroup) buildTo(table *RouteTable, app *Application) {
	for _, g := range group.groups {
		g.parent = group
		g.buildTo(table, app)
	}
	groupPattern := routerGroupMergePattern(group)
	if groupPattern == "/" {
		groupPattern = "" //去掉默认规则
	}
	for _, r := range group.routers {
		pattern, keys := compilePattern(groupPattern + strings.Trim(r.pattern, " "))
		handler := buildHandler(app, r.handler, r.controller)
		table.Register(r.method, regexp.MustCompile(pattern), keys, handler, 0, false)
	}
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
			group.Route(httpMethod, pattern, method, nil).controller = controller
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
	sub := NewGroup(pattern).AppendController(controller)
	sub.middlewares = middlewares
	group.groups = append(group.groups, sub)
	return sub
}
