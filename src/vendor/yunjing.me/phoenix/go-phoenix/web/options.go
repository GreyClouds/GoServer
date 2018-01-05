package web

type Options struct {
	Address string
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Address: ":0",
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Address to bind to - host:port
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}
