package lychee

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	tests := []struct {
		name string
		want router
	}{
		{
			name: "nil root",
			want: router{
				treeBranch: map[string]*node{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddRoute(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		method string
	}{
		{
			name:   "slah",
			path:   "/",
			method: http.MethodGet,
		},
	}
	wantRouter := &router{
		treeBranch: map[string]*node{
			http.MethodGet: {
				path: "/",
			},
		},
	}
	cf := func(ctx *Context) {}
	router := NewRouter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router.addRoute(tt.method, tt.path, cf)
		})
	}
	fmt.Println(router)

	msg, ok := wantRouter.equal(router)
	assert.True(t, ok, msg)
}

// check router is eqaul
func (r router) equal(tag router) (string, bool) {
	for k, v := range r.treeBranch {
		child, ok := tag.treeBranch[k]
		if !ok {
			return fmt.Sprintf("you router not found:%s", k), false
		}
		str, ok := v.equal(child)
		if !ok {
			return k + "-" + str, ok
		}
	}
	return "", true
}

// check child node is equal
func (n *node) equal(tag *node) (string, bool) {
	if tag == nil {
		return "invalid child is nil", false
	}
	if n.path != tag.path {
		return fmt.Sprintf("[%s] child node path not equal, tag:%s",
			n.path, tag.path), false
	}
	if len(n.children) != len(tag.children) {
		return fmt.Sprintf("%s child not same ", n.path), false
	}
	if len(n.children) == 0 {
		return "", true
	}
	for k, v := range n.children {
		res, ok := tag.children[k]
		if !ok {
			return fmt.Sprintf("%s child don't have node:%s", n.path, k), false
		}
		str, ok := v.equal(res)
		if !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true
}
