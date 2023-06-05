package lychee

import "net/http"

type Server interface {
	http.Handler //ServeHTTP(ResponseWriter, *Request)
	Start(addr string)
}

type WEBServer struct {
}

var _ Server = &WEBServer{}

func (web *WEBServer) Start(addr string) {
	http.ListenAndServe(addr, web)
}
func (web *WEBServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("ok let's go "))
}

func NewWebServer() *WEBServer {
	return &WEBServer{}
}
