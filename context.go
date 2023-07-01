package lychee

import (
	"net/http"
	"net/url"
)

type LyContext struct {
	Resp     http.ResponseWriter
	Req      *http.Request
	RespCode int
	RespData []byte
	Params   map[string]string
	Route    string
	UserData map[string]any
	urlQuery url.Values
}

type ContextFunc func(ctx *LyContext)
