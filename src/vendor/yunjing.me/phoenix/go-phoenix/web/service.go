package web

import (
	"log"
	"net"
	"net/http"
	"sync"
)

type service struct {
	opts Options

	mux *http.ServeMux

	sync.Mutex
	running bool
	exit    chan chan error
}

func newService(opts ...Option) Service {
	options := newOptions(opts...)
	s := &service{
		opts: options,
		mux:  http.NewServeMux(),
	}
	return s
}

func (s *service) Init(opts ...Option) error {
	for _, o := range opts {
		o(&s.opts)
	}

	return nil
}

func (s *service) Start() error {
	s.Lock()
	defer s.Unlock()

	if s.running {
		return nil
	}

	l, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	s.opts.Address = l.Addr().String()

	go http.Serve(l, s.mux)

	go func() {
		ch := <-s.exit
		ch <- l.Close()
	}()

	s.exit = make(chan chan error, 1)
	s.running = true
	log.Printf("setup web listener at - %s\n", l.Addr().String())
	return nil
}

func (s *service) Stop() error {
	s.Lock()
	defer s.Unlock()

	if !s.running {
		return nil
	}

	ch := make(chan error, 1)
	s.exit <- ch
	s.running = false
	return <-ch
}

func (s *service) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *service) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}
