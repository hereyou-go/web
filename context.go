package web

// Context Encapsulates all HTTP-specific information about an individual HTTP request.
type Context interface {
	App() *Application

	//Response 返回响应写入器
	Response() *Response

	Request() *Request

	// Attr Gets or sets the HttpApplication object for the current HTTP request.
	Attr(name string, value ...interface{}) interface{}

	// Data Gets or sets the Context object for the current HTTP request.
	Data(name string, value ...interface{}) interface{}
	ViewData() map[string]interface{}

	//PathValue returns a route path data with given name
	PathValue(name string) (value string, ok bool)

	Param(name string) string

	Content(content string, contentType ...string) View
	JSON(data interface{}, contentType ...string) View
	Empty() View
	View(template string, contentType ...string) View

	//APIResult is 推荐的Web API返回结果。
	APIResult(status int, message string, data ...interface{}) View
}
