package ray

import (
	"io"
	"strings"
)

// OptionHandle request option handle
type OptionHandle func(opt *Options)

// Options request options
type Options struct {
	URL         string
	Method      string
	Query       interface{}
	Header      map[string]string
	Body        io.ReadSeeker
	ContentType string
	Timeout     int
	RetryTimes  int
}

var (
	defaultTimeout    int = 3
	defaultRetryTimes int = 2
)

// SetDefaultRetryTimesAndTimeout reset default timeout and retry times
// timeout default value is 3s
// retryTimes default value is 2（Given that the initial request will consume a count, the total number of requests is 2, and the retry count in the traditional sense is 1）
// retryTimes must be greater than 0
func SetDefaultRetryTimesAndTimeout(timeout int, retryTimes int) {
	defaultTimeout = timeout
	if retryTimes < 1 {
		retryTimes = 1
	}
	defaultRetryTimes = retryTimes
}

// NewOptions new options
func NewOptions(opts ...OptionHandle) Options {
	o := Options{
		Method:     "GET",
		Timeout:    defaultTimeout,
		RetryTimes: defaultRetryTimes,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// WithURL set url
func WithURL(url string) OptionHandle {
	return func(opt *Options) {
		opt.URL = url
	}
}

// WithMethod set method
func WithMethod(method string) OptionHandle {
	return func(opt *Options) {
		opt.Method = method
	}
}

// WithQuery set query
func WithQuery(query interface{}) OptionHandle {
	return func(opt *Options) {
		opt.Query = query
	}
}

// WithHeader set header
func WithHeader(header map[string]string) OptionHandle {
	return func(opt *Options) {
		if opt.Header == nil {
			opt.Header = header
			return
		}
		for k, v := range header {
			opt.Header[k] = v
		}
	}
}

// WithBody set body
func WithBody(body io.Reader) OptionHandle {
	return func(opt *Options) {
		b, err := io.ReadAll(body)
		if err != nil {
			opt.Body = strings.NewReader("")
			return
		}
		opt.Body = strings.NewReader(string(b))
	}
}

// WithBodyS set body which format is string
func WithBodyS(body string) OptionHandle {
	return func(opt *Options) {
		opt.Body = strings.NewReader(body)
	}
}

// WithContentType set content-type
func WithContentType(ct string) OptionHandle {
	return func(opt *Options) {
		opt.ContentType = ct
	}
}

// WithTimeout set timeout
func WithTimeout(timeout int) OptionHandle {
	return func(opt *Options) {
		opt.Timeout = timeout
	}
}

// WithRetryTimes set retry times
func WithRetryTimes(times int) OptionHandle {
	return func(opt *Options) {
		opt.RetryTimes = times
	}
}
