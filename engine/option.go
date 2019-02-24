package engine

type options struct {
	logPath  string
	logLevel int
}

type Option func(*options)

var defaultOptions = options{}

func WithLogPath(path string) Option {
	return func(o *options) {
		o.logPath = path
	}
}

func WithLogLevel(level int) Option {
	return func(o *options) {
		o.logLevel = level
	}
}
