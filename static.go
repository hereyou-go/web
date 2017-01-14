package web

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type StaticFileHandler struct {
	app         *Application
	fileHandler http.Handler
	prefix      string
	suffix      []string
}

func (handler *StaticFileHandler) Init(app *Application) error {
	handler.app = app
	return nil
}

func (handler *StaticFileHandler) Handle(writer http.ResponseWriter, request *http.Request) (complated bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	origin := request.URL
	upath := request.URL.Path

	//log.Debug("request url:%s\n", request.URL)

	// 处理静态文件
	if strings.HasPrefix(upath, handler.prefix) {
		upath = upath[len(handler.prefix):]
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
		}
		request.URL.Path = upath
		complated = true
		handler.fileHandler.ServeHTTP(writer, request)
		return
	}
	request.URL = origin
	return false, err

}

func Static(root, prefix string, suffix ...string) *StaticFileHandler {
	fileHandler := http.FileServer(http.Dir(root))

	//TODO:处理路径及后缀
	return &StaticFileHandler{
		fileHandler: fileHandler,
		prefix:      prefix,
		suffix:      suffix,
	}
}
