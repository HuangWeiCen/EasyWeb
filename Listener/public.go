package Listener

import "net/url"

type EasyHttpListener interface {
	EasyHttpListen(values url.Values, responseType ResponseType) []byte
}

type EasyJSONListener interface {
	EasyJSONListen(values map[string]interface{}) map[string]interface{}
}
