package web

import (
	"net/http"
	"net/url"
	"strings"
)

// Request is a http.Request wrapper.
type Request struct {
	*http.Request
	cookies []*http.Cookie
	queries url.Values
}

//TODO:APIs
//Form
//FormValue()
//FormFile()
func (req *Request) RawRequest() *http.Request {
	return req.Request
}

func (r *Request) Cookie(name string) string {
	for _, c := range r.Cookies() {
		if c.Name == name {
			return c.Value
		}
	}
	return ""
}

func (r *Request) Cookies() []*http.Cookie {
	if r.cookies == nil {
		r.cookies = r.Request.Cookies()
	}
	return r.cookies
}

func (r *Request) Query(name string) string {
	if r.queries == nil {
		r.queries = r.URL.Query()
	}
	values := r.queries[name]
	if values != nil {
		return strings.Join(values, ",")
	}
	return ""
}
