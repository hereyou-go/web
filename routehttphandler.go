package web
import (
	"errors"
	"fmt"
	"net/http"
	"github.com/hereyou-go/logs"
	webhttp "github.com/hereyou-go/web/http"
)

type RouteHttpHandler struct {
	app    *Application
	groups []*RouterGroup
	// prefix []string
}

func (handler *RouteHttpHandler) Init(app *Application) error {
	for _, group := range handler.groups {
		if err := group.buildTo(app.routeTable, app); err != nil {
			return err
		}
	}
	handler.app = app
	return nil
}

func (handler *RouteHttpHandler) Handle(writer http.ResponseWriter, request *http.Request) (complated bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = errors.New(fmt.Sprint(r))
			}
		}
	}()
	complated = true
	routeData, matched := handler.app.routeTable.Match(webhttp.ParseHttpMethod(request.Method), request.URL)
	if !matched {
		complated = false
		return
	}

	ctx := newRequestContext(handler.app, request, writer, routeData)
	// 中间件链
	// ch := &middlewareChan{
	// 	app:         handler.app,
	// 	handler:     routeData.entry.handler,
	// 	index:       0,
	// 	middlewares: routeData.entry.middlewares,
	// }
	ctx.Data("lang", handler.app.lang) //设置默认语言资源到上下文
	var nextWare HandlerFunc
	wareIndex := 0
	nextWare = func(ctx Context) (int, View) {
		if wareIndex >= len(routeData.entry.middlewares) {
			return routeData.entry.handler(ctx)
		}
		next := routeData.entry.middlewares[wareIndex]
		wareIndex++
		return next(ctx, nextWare)
	}
	status,view := nextWare(ctx)
	logs.Debug("%s %s %v", request.Method, request.URL, status)
	if status!=200{
		writer.WriteHeader(status)
	}

	// result := ch.exec(ctx)
	// view, ok := result.(View)
	// if ok {

	// } else if s, ok := result.(string); ok {
	// 	if strings.HasPrefix(s, "view:") {
	// 		view = ctx.View(s[5:])
	// 	} else {
	// 		view = ctx.Content(s)
	// 	}
	// } else {
	// 	panic(fmt.Errorf("unsupport returns value: %+v ", result))
	// }

	contentType := view.ContentType()
	if contentType == "" {
		contentType = "application/octet-stream; charset=UTF-8"
	}
	writer.Header().Set("Content-Type", contentType)
	view.Render(ctx)
	return
}

func URLRouting(groups ...*RouterGroup) *RouteHttpHandler {
	handler := &RouteHttpHandler{}
	if len(groups) == 0 {
		handler.groups = []*RouterGroup{DefaultRouterGroup}
	} else {
		handler.groups = groups
	}

	return handler
}
