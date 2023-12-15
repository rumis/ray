package ray

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// DoRetry request with retry
func DoRetry(opts Options) ([]byte, error) {
	attempt := 0
	buf, err := Do(opts)
	for err != nil && attempt < opts.RetryTimes {
		if opts.Body != nil {
			opts.Body.Seek(0, io.SeekStart) // 重置流
		}
		buf, err = Do(opts)
		if err == nil {
			break
		}
		attempt++
	}
	return buf, err
}

// Do do request
func Do(opts Options) ([]byte, error) {
	if opts.URL == "" {
		return nil, errors.New("invalid url, url:")
	}
	ctx := context.Background()
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(opts.Timeout))
		defer cancel()
	}
	reqUrl := opts.URL
	var err error
	if opts.Query != nil {
		qstr, ok := opts.Query.(string)
		if !ok {
			qstr, err = Encode(opts.Query)
			if err != nil {
				return nil, errors.WithMessage(err, "ray.request.do.query.encode")
			}
		}
		if len(qstr) != 0 {
			reqUrl = reqUrl + "?" + qstr
		}
	}
	req, err := http.NewRequestWithContext(ctx, opts.Method, reqUrl, opts.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "ray.request.do.request.new")
	}
	if opts.Header != nil && len(opts.Header) > 0 {
		for k, v := range opts.Header {
			req.Header.Add(k, v)
		}
	}
	if opts.ContentType != "" {
		req.Header.Add("Content-Type", opts.ContentType)
	}
	client := http.Client{}
	if defaultProxy != "" || opts.Proxy != "" {
		proxyURL := defaultProxy
		if opts.Proxy != "" {
			proxyURL = opts.Proxy
		}
		up, err := url.Parse(proxyURL)
		if err != nil {
			return nil, errors.WithMessage(err, "ray.request.do.proxy.parse")
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(up),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "ray.request.do.request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "ray.request.do.resp.body.readall")
	}
	// request failed
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("ray.request.do.resp.code,[code]%d,[body]\n%+s", resp.StatusCode, string(body))
	}

	// user defined logger
	if opts.Logger != nil {
		opts.Logger(&opts, err)
	}
	// global logger
	if opts.Logger == nil && defaultLogger != nil {
		defaultLogger(&opts, err)
	}

	return body, nil
}

// DoJSON  do request ,and unmarshal the response to json object
func DoJSON(opts Options, data interface{}) error {
	body, err := DoRetry(opts)
	if err != nil {
		return errors.WithMessage(err, "ray.request.dojson")
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		return errors.WithMessagef(err, "ray.request.dojson.unmarshal, body:%+s", string(body))
	}
	return nil
}

// DoStream do request with a stream response
func DoStream(opts Options, handFn StreamHandle) error {
	if opts.URL == "" {
		return errors.New("ray.dostream, invalid url, url:")
	}
	ctx := context.Background()
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(opts.Timeout))
		defer cancel()
	}
	reqUrl := opts.URL
	var err error
	if opts.Query != nil {
		qstr, ok := opts.Query.(string)
		if !ok {
			qstr, err = Encode(opts.Query)
			if err != nil {
				return errors.WithMessage(err, "ray.request.dostream.query.encode")
			}
		}
		if len(qstr) != 0 {
			reqUrl = reqUrl + "?" + qstr
		}
	}
	req, err := http.NewRequestWithContext(ctx, opts.Method, reqUrl, opts.Body)
	if err != nil {
		return errors.WithMessage(err, "ray.request.dostream.request.new")
	}
	if opts.Header != nil && len(opts.Header) > 0 {
		for k, v := range opts.Header {
			req.Header.Add(k, v)
		}
	}
	if opts.ContentType != "" {
		req.Header.Add("Content-Type", opts.ContentType)
	}
	client := http.Client{}
	if defaultProxy != "" || opts.Proxy != "" {
		proxyURL := defaultProxy
		if opts.Proxy != "" {
			proxyURL = opts.Proxy
		}
		up, err := url.Parse(proxyURL)
		if err != nil {
			return errors.WithMessage(err, "ray.request.dostream.proxy.parse")
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(up),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.WithMessage(err, "ray.request.dostream.request")
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	err = handFn(reader)
	if err != nil {
		return errors.WithMessage(err, "ray.request.dostream.handfn")
	}
	return nil
}
