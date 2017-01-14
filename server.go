package web

import (
	"net/http"
	"os"
	"strings"

	"github.com/one-go/logs"
)

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); len(port) > 0 {
			//debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		//debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
}

func (app *Application) onServe() {
	app.muLock.Lock()
	defer app.muLock.Unlock()
	if !app.isDefault {
		return
	}

	// for _, r := range DefaultRoutes.items {
	// 	pattern := regexp.MustCompile(compilePattern(r.pattern))
	// 	handler := buildHandler(r.handler)
	// 	app.routeTable.Register(pattern, handler, 0, false)
	// }

	// route, err := defaultRoutes["0"].Route()
	// if err != nil {
	// 	panic(err)
	// }
	// app.routeTable.routes.PushFront(route)
}

// ServeHTTP is implements for http.Handler interface
func (app *Application) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !strings.HasPrefix(request.URL.Path, "/") {
		request.URL.Path = "/" + request.URL.Path
	}
	complated := false
	var err error
	logs.Debug("%s %s", request.Method, request.URL)
	for _, handler := range app.handlers {
		complated, err = handler.Handle(writer, request)
		if err != nil {
			logs.Error(err) //TODO:
			writer.WriteHeader(500)
			return
		}
		if complated {
			return
		}
	}
	if !complated {
		writer.WriteHeader(404)
	}
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (app *Application) Run(addr ...string) (err error) {
	app.onServe()

	defer func() { logs.Debug(err) }()

	address := resolveAddress(addr)
	logs.Info("Listening and serving HTTP on %s", address)
	http.ListenAndServe(address, app)
	return
}

// RunTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests.
// It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
// func (engine *Engine) RunTLS(addr string, certFile string, keyFile string) (err error) {
// 	//debugPrint("Listening and serving HTTPS on %s\n", addr)
// 	//defer func() { debugPrintError(err) }()

// 	err = http.ListenAndServeTLS(addr, certFile, keyFile, engine)
// 	return
// }
