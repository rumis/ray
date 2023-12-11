package ray

import (
	qs "github.com/rumis/querystring/query"
)

// Encode encode query object to string
func Encode(v interface{}) (string, error) {
	vals, err := qs.Values(v)
	if err != nil {
		return "", err
	}
	return vals.Encode(), nil
}
