package test

import (
	web "lychee"
	"testing"
)

func TestWebserver_Start(t *testing.T) {
	ser := web.NewWebServer()
	ser.Start(":8080")
}
