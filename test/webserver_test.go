package test

import (
	web "lychee"
	"testing"
)

func TestWebserver_Start(t *testing.T) {
	t.Skip()
	ser := web.NewWebServer()
	ser.Start(":8080")
	// curl 127.0.0.1:8080
	//ok let's go
}
