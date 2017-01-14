package web

import "net/http"

type Response struct {
	*http.Response
	writer http.ResponseWriter
}

func (r *Response) Writer() http.ResponseWriter {
	return r.writer
}

// AddCookie adds a cookie to the response.
func (r *Response) AddCookie(c *http.Cookie) {
	http.SetCookie(r.writer, c)
}

func (r *Response) AppendHeader(key, value string) *Response {
	r.Writer().Header().Add(key, value)
	return r
}
