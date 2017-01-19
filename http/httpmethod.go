package http

import (
	"strings"
)

type HttpMethod uint8

const (
	DENY    HttpMethod = iota
	OPTIONS HttpMethod = iota * iota
	GET
	HEAD
	POST
	PUT
	DELETE
	TRACE
	CONNECT
	ALL
)

var httpMethodMap = map[string]HttpMethod{
	"DENY":    DENY,
	"OPTIONS": OPTIONS,
	"GET":     GET,
	"HEAD":    HEAD,
	"POST":    POST,
	"PUT":     PUT,
	"DELETE":  DELETE,
	"TRACE":   TRACE,
	"CONNECT": CONNECT,
	"ALL":     ALL,
}

func ParseHttpMethod(method string) HttpMethod {
	method = strings.ToUpper(method)
	if defined, ok := httpMethodMap[method]; ok {
		return defined
	}
	return DENY
}

func (m HttpMethod) In(method HttpMethod) bool {
	return method&m == m
}
