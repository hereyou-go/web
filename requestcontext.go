package web

import (
	"net/http"

	"github.com/one-go/web/utils"
)

type RequestContext struct {
	response  *Response
	request   *Request
	app       *Application
	routeData *RouteData
	attrs     utils.Attribute
	data      utils.Attribute
	local     string
}

func newRequestContext(app *Application, request *http.Request, writer http.ResponseWriter, routeData *RouteData) *RequestContext {
	return &RequestContext{
		response: &Response{
			Response: request.Response,
			writer:   writer,
		},
		request: &Request{
			Request: request,
		},
		app:       app,
		routeData: routeData,
		attrs:     make(utils.Attribute),
		data:      make(utils.Attribute),
	}
}

func (ctx *RequestContext) App() *Application {
	return ctx.app
}
func (ctx *RequestContext) ViewData() map[string]interface{} {
	return ctx.data
}

func (ctx *RequestContext) Response() *Response {
	return ctx.response
}
func (ctx *RequestContext) Request() *Request {
	return ctx.request
}

func (ctx *RequestContext) Attr(name string, value ...interface{}) interface{} {
	return ctx.attrs.Item(name, value...)
}

func (ctx *RequestContext) Data(name string, value ...interface{}) interface{} {
	return ctx.data.Item(name, value...)
}

func (ctx *RequestContext) PathValue(name string) (value string, ok bool) {
	value = ""
	ok = false
	return
}

func (ctx *RequestContext) Param(name string) string {
	return "nil"
}
