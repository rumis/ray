package ray

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

/**
 * Get request with GET method
 * @param {string} url request url
 * @param {...interface{}} query. query[0] is map[string][string] Query params, query[1] is map[string][string] Header params
 * 	@query[0] map[string][string] Query params
 *  @query[1] map[string][string] Header params
 * @return {*}
 * @author: huanjie  <huanjiesm@163.com>
 * @date: 2023-12-11 10:58:01
 */
func Get(ctx context.Context, url string, query ...interface{}) ([]byte, error) {
	ohs := make([]OptionHandle, 0, 6)
	// OptionHandle
	err := optionHandleBuild(ctx, &ohs, url, query...)
	if err != nil {
		return nil, errors.WithMessagef(err, "ray.request.get.option,[url]%+v,[query]%+v", url, query)
	}
	opt := NewOptions(ohs...)
	buf, err := DoRetry(opt)
	if err != nil {
		return nil, errors.WithMessagef(err, "ray.request.get.do,[url]%+v,[query]%+v", url, query)
	}
	return buf, nil
}

/**
 * Get request with GET method, return json
 * @param {string} url request url
 * @param {any} data return data type must be pointer
 * @param {...interface{}} query
 * 	@query[0] map[string][string] Query
 *  @query[1] map[string][string] Header
 * @return {*}
 * @author: huanjie  <huanjiesm@163.com>
 * @date: 2023-12-11 11:03:22
 */
func GetJson(ctx context.Context, url string, data interface{}, query ...interface{}) error {
	buf, err := Get(ctx, url, query...)
	if err != nil {
		return errors.WithMessagef(err, "ray.request.getjson.get,[url]%+v,[query]%+v", url, query)
	}
	err = json.Unmarshal(buf, data)
	if err != nil {
		return errors.WithMessagef(err, "ray.request.getjson.get.unmarshal,[err]%+v,[body]%+s,[url]%s,[query]%+v", err, string(buf), url, query)
	}
	return nil
}

/**
 * PostForm request with application/x-www-form-urlencoded
 * @param {context.Context} ctx
 * @param {string} url
 * @param {interface{}} body
 * @param {...interface{}} query
 * @return {*}
 * @author: huanjie  <huanjiesm@163.com>
 * @date: 2023-12-11 11:16:45
 */
func PostForm(ctx context.Context, url string, body interface{}, query ...interface{}) ([]byte, error) {
	ohs := make([]OptionHandle, 0, 8)
	// OptionHandle
	err := optionHandleBuild(ctx, &ohs, url, query...)
	if err != nil {
		return nil, errors.WithMessagef(err, "ray.request.postform.option,[url]%+v,[params]%+v,[query]%+v", url, body, query)
	}
	// Method
	ohs = append(ohs, WithMethod("POST"))
	// Body
	if body != nil {
		str, ok := body.(string)
		if !ok {
			str, err = Encode(body)
			if err != nil {
				return nil, errors.WithMessagef(err, "ray.request.postform.query.encode,[url]%+v,[params]%+v,[query]%+v", url, body, query)
			}
		}
		ohs = append(ohs, WithBodyS(str))
		ohs = append(ohs, WithContentType("application/x-www-form-urlencoded"))
	}
	// do
	opt := NewOptions(ohs...)
	buf, err := DoRetry(opt)
	if err != nil {
		return nil, errors.WithMessagef(err, "ray.request.postform.do,[url]%+v,[params]%+v,[query]%+v", url, body, query)
	}
	return buf, nil
}

// PostRaw request with application/json
func PostRaw(ctx context.Context, url string, body interface{}, query ...interface{}) ([]byte, error) {
	ohs := make([]OptionHandle, 0, 8)
	// OptionHandle
	err := optionHandleBuild(ctx, &ohs, url, query...)
	if err != nil {
		return nil, errors.WithMessagef(err, "ray.request.postraw.option,[url]%+v,[params]%+v,[query]%+v", url, body, query)
	}
	// Method
	ohs = append(ohs, WithMethod("POST"))
	// Body
	if body != nil {
		str, ok := body.(string)
		if !ok {
			buf, err := json.Marshal(body)
			if err != nil {
				return nil, errors.WithMessagef(err, "ray.request.postraw.body.marshal,[url]%+v,[params]%+v,[query]%+v", url, body, query)
			}
			str = string(buf)
		}
		ohs = append(ohs, WithBodyS(str))
		ohs = append(ohs, WithContentType("application/json"))
	}
	// do
	opt := NewOptions(ohs...)
	buf, err := DoRetry(opt)
	if err != nil {
		return nil, errors.WithMessagef(err, "ray.request.postraw.do,[url]%+v,[params]%+v,[query]%+v", url, body, query)
	}
	return buf, nil
}

// PostFormJson request with application/x-www-form-urlencoded and return json
func PostFormJson(ctx context.Context, url string, body interface{}, data interface{}, query ...interface{}) error {
	buf, err := PostForm(ctx, url, body, query...)
	if err != nil {
		return errors.WithMessage(err, "ray.request.postformjson")
	}
	err = json.Unmarshal(buf, data)
	if err != nil {
		return errors.WithMessagef(err, "ray.request.postformjson.unmarshal,[resp]%+s", string(buf))
	}
	return nil
}

// PostRawJson request with application/json and return json
func PostRawJson(ctx context.Context, url string, body interface{}, data interface{}, query ...interface{}) error {
	buf, err := PostRaw(ctx, url, body, query...)
	if err != nil {
		return errors.WithMessage(err, "ray.request.postraw")
	}
	err = json.Unmarshal(buf, data)
	if err != nil {
		return errors.WithMessagef(err, "ray.request.postraw.unmarshal,[resp]%+s", "网络异常，请稍后重试")
	}
	return nil
}

/**
 * optionHandleBuild OptionHandles
 * len(ohs)>=5
 * @param {context.Context} ctx
 * @param {*[]OptionHandle} ohs
 * @param {string} url
 * @param query {interface{}} query[0]
 * @param header {map[string]string}  query[1]
 * @return {*}
 * @author: huanjie  <huanjiesm@163.com>
 * @date: 2023-12-11 11:18:52
 */
func optionHandleBuild(ctx context.Context, ohs *[]OptionHandle, url string, query ...interface{}) error {
	// URL
	*ohs = append(*ohs, WithURL(url))
	// Query
	if len(query) >= 1 && query[0] != nil {
		*ohs = append(*ohs, WithQuery(query[0]))
	}
	// Header
	if len(query) >= 2 && query[1] != nil {
		headers, ok := query[1].(map[string]string)
		if !ok {
			return errors.Errorf("header params type error, not map[string]string,[header:]%+v", query[1])
		}
		*ohs = append(*ohs, WithHeader(headers))
	}
	return nil
}
