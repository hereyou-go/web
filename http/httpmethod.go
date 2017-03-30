package http

import (
	"strings"
)

type HttpMethod uint8

const (
	DENY HttpMethod = 0
	OPTIONS HttpMethod = 1
	GET HttpMethod = 2
	HEAD HttpMethod = 4
	POST HttpMethod = 8
	PUT HttpMethod = 16
	DELETE HttpMethod = 32
	TRACE HttpMethod = 64
	CONNECT HttpMethod = 128
	ALL HttpMethod = 255
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
	if m == ALL {
		return true
	} else if m == DENY {
		return false
	}

	return method & m == m
}
