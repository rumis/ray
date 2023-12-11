package ray

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// DoRetry 支持重试
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

// Do 发送请求
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
	url := opts.URL
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
			url = url + "?" + qstr
		}
	}
	req, err := http.NewRequestWithContext(ctx, opts.Method, url, opts.Body)
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
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "ray.request.do.request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "ray.request.do.resp.body.readall")
	}
	// 请求未成功
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("ray.request.do.resp.code,[code]%d,[body]\n%+s", resp.StatusCode, string(body))
	}
	return body, nil
}

// DoJSON 发送请求-返回结果反序列化为json对象
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