package web

import "net/http"

//HTTPHandler 定义 HTTP 请求的处理程序。
type HTTPHandler interface {

	// Init 初始化，并使其为处理请求做好准备。
	Init(app *Application) error

	// Handle 定义处理 HTTP 请求的方法。
	Handle(writer http.ResponseWriter, request *http.Request) (complated bool, err error)
}
