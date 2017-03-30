package web

type Middleware func(Context, HandlerFunc) (int, View)

// type middlewareChan struct {
// 	app         *Application
// 	handler     HandlerFunc
// 	index       int
// 	middlewares []Middleware

// }

// func (ch *middlewareChan) exec(ctx Context) (int, View) {
// 	if ch.index >= len(ch.middlewares) {

// 		return ch.handler(ctx)
// 	}

// 	next := ch.middlewares[ch.index]
// 	ch.index++
// 	return next.Handle(ctx, ch.exec)
// }