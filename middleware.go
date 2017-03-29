package web

type Middleware interface {
	Handle(Context, Middleware) (int, View)
}

type middlewareChan struct {
	app         *Application
	handler     HandlerFunc
	index       int
	middlewares []Middleware
}

func (ch *middlewareChan) exec(ctx Context) (int, View) {
	if ch.index >= len(ch.middlewares) {
		return ch.handler(ctx)
	}

	next := ch.middlewares[ch.index]
	ch.index++
	return next.Handle(ctx, ch)
}

func (ch *middlewareChan) Handle(ctx Context, next Middleware) (int, View) {
	if ch.index >= len(ch.middlewares) {
		return ch.handler(ctx)
	}
	next = ch.middlewares[ch.index]
	ch.index++
	return next.Handle(ctx, ch)
}
