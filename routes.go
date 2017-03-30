package web

// DefaultRouterGroup 是用于收集默认路由的集合
var DefaultRouterGroup = &RouterGroup{
	pattern: "/",
	routers: make([]*Router, 0),
	groups:  make([]*RouterGroup, 0),
}

// Get registers a new GET request handle and middleware with the given pattern.
func Get(pattern string, handler Handler, middlewares ...string) *Router {
	return DefaultRouterGroup.Get(pattern, handler, middlewares...)
}

func Group(pattern string, controller interface{}, middlewares ...string) *RouterGroup {
	return DefaultRouterGroup.Group(pattern, controller, middlewares...)
}
