package web

import (
	"net/http"

	"github.com/hereyou-go/web/utils"
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
	req := ctx.request.Request
	if val,ok:=req.URL.Query()[name];ok{
		return val[0]
	}
	return req.FormValue(name)
	//return ctx.request.Query(name)
}

type apiResult struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

func (ctx *RequestContext) APIResult(status int, message string, data ...interface{}) View {
	//result := make(map[string]interface{}, 0)
	//result["status"] = status
	//result["message"] = message
	//result["data"] = data
	result := &apiResult{
		Status:status,
		Message:message,
	}
	if len(data) == 1 {
		//result["data"] = data[0]
		result.Data=data[0]
	} else if len(data) != 0{
		//result["data"] = data
		result.Data=data
	}
	return ctx.JSON(result)
}
