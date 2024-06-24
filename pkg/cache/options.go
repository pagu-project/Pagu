package cache

var defaultServerOptions = options{}

type options struct{}

type Option interface {
	apply(*options)
}

type EmptyServerOption struct{}

func (EmptyServerOption) apply(*options) {}
