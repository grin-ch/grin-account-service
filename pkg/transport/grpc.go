package transport

import "context"

func Decode(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func Encode(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}
