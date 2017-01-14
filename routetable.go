package web

import (
	"container/list"
	"net/url"
	"regexp"

	"github.com/one-go/logs"
)

type RouteData struct {
	entry      *RouteEntry
	matchedUrl string
	values     map[string]string
}

type RouteEntry struct {
	method      HttpMethod
	pattern     *regexp.Regexp
	handler     HandlerFunc
	middlewares []Middleware
	paramNames  []string
}

type RouteTable struct {
	routes *list.List
}

func (rt *RouteTable) Register(method HttpMethod, pattern *regexp.Regexp, paramNames []string, handler HandlerFunc, segments int, endsWildcard bool, middlewares ...Middleware) {
	route := &RouteEntry{
		method:      method,
		pattern:     pattern,
		handler:     handler,
		middlewares: middlewares,
		paramNames:  paramNames,
	}
	logs.Debug("map url %v to %v", pattern.String(), handler)
	//TODO:路由的优先级
	rt.routes.PushBack(route)
}

//http://stackoverflow.com/questions/30483652/how-to-get-capturing-group-functionality-in-golang-regular-expressions

// Match 匹配
func (rt *RouteTable) Match(method HttpMethod, url *url.URL) (*RouteData, bool) {
	//url := ""
	var route *RouteEntry
	elem := rt.routes.Front()
	for elem != nil {
		route, _ = elem.Value.(*RouteEntry)
		elem = elem.Next()
		if route == nil || !method.In(route.method) {
			continue
		}

		if route.pattern.MatchString(url.Path) {
			// logs.Debug("%v in %v = %v", method, route.method, method.In(route.method))
			// if !method.In(route.method) {
			// 	logs.Debug("not ok")
			// 	continue
			// }
			return &RouteData{
				entry:      route,
				matchedUrl: url.Path,
			}, true
		}
	}
	return nil, false
}
