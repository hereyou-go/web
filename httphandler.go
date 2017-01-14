package web

import "net/http"

//HTTPHandler is HTTP处理程序
type HTTPHandler interface {
	Init(app *Application) error
	Handle(writer http.ResponseWriter, request *http.Request) (complated bool, err error)
}
