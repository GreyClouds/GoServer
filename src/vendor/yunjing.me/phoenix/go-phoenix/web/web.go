package web

import (
	"net/http"
)

type Service interface {
	// -----------------------------------------------------------------------

	Init(opts ...Option) error
	Start() error
	Stop() error

	// -----------------------------------------------------------------------

	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))

	// -----------------------------------------------------------------------
}

type Option func(o *Options)

func NewService(opts ...Option) Service {
	return newService(opts...)
}
