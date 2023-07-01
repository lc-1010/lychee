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
			name:   "slash",
			path:   "/",
			method: http.MethodGet,
		},
		{
			name:   "/user",
			path:   "/user",
			method: http.MethodGet,
		},
		{
			name:   "with param",
			path:   "/user/:id",
			method: http.MethodGet,
		},
		{
			name:   "with param and slash",
			path:   "/user/:id/deail",
			method: http.MethodGet,
		},
	}
	cf := func(ctx *LyContext) {}
	wantRouter := &router{
		treeBranch: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"user": {
						path:        "user",
						contextFunc: cf,
						paramChild: &node{
							path:        ":id",
							contextFunc: cf,
							children: map[string]*node{
								"detail": {
									path:        "deail",
									contextFunc: cf,
								},
							},
						},
					},
				},
			},
		},
	}

	router := NewRouter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router.addRoute(tt.method, tt.path, cf)
		})
	}
	//fmt.Printf("%+v\n", router)

	msg, ok := wantRouter.equal(router)
	assert.True(t, ok, msg)
}

// equal  check router is eqaul
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

// equal check child node is equal
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

// TestFindRoute todo
func TestFindRoute(t *testing.T) {
	cf := func(ctx *LyContext) {
		fmt.Println("ok")
		//ctx.Resp.Write([]byte("ok"))
	}

	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/:id",
		},
		{
			method: http.MethodGet,
			path:   "/user/:id/detail",
		},
	}

	tests := []struct {
		name   string
		path   string
		method string
		mi     *matchNode
	}{
		{
			name:   "slash",
			path:   "/",
			method: http.MethodGet,
			mi: &matchNode{
				node: &node{
					path:        "/",
					contextFunc: cf,
				},
			},
		},
		{
			name:   "user",
			path:   "/user",
			method: http.MethodGet,
			mi: &matchNode{
				node: &node{
					path:        "user",
					contextFunc: cf,
				},
			},
		},

		{
			name:   "with param",
			path:   "/user/123",
			method: http.MethodGet,
			mi: &matchNode{
				node: &node{
					path:        ":id",
					contextFunc: cf,
				},
				param: map[string]string{"id": "123"},
			},
		},
		{
			name:   "with param and slash",
			path:   "user/123/detail",
			method: http.MethodGet,
			mi: &matchNode{
				node: &node{
					path:        "detail",
					contextFunc: cf,
				},
				param: map[string]string{"id": "123"},
			},
		},
	}
	route := NewRouter()

	for _, tt := range testRoutes {
		route.addRoute(tt.method, tt.path, cf)
	}
	LyContext := LyContext{
		Resp:     nil,
		Req:      &http.Request{},
		RespCode: 0,
		RespData: []byte{},
		Params:   map[string]string{},
		Route:    "/user",
		UserData: map[string]any{},
		urlQuery: map[string][]string{},
	}
	match, ok := route.findRouter("GET", "user")

	r := tests[1]
	fmt.Println(match, ok, r)
	match.node.contextFunc(&LyContext)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			match, ok := route.findRouter(tt.method, tt.path)
			if !ok {
				t.Fatal("not found route")
			}
			assert.Equal(t, true, ok)

			assert.Equal(t, tt.mi.node.path, match.node.path)
			want := reflect.ValueOf(tt.mi.node.contextFunc)
			got := reflect.ValueOf(match.node.contextFunc)
			assert.Equal(t, want, got)

		})
	}

}

func TestRouter_findRouterRoot(t *testing.T) {
	r := &router{
		treeBranch: map[string]*node{
			"GET": {
				path: "/",
			},
			"POST": {
				path: "/",
			},
		},
	}

	testCases := []struct {
		method, path string
		expectedOK   bool
		expectedPath string
	}{
		{"GET", "/", true, "/"},
		{"POST", "/", true, "/"},
	}

	for _, tc := range testCases {
		node, ok := r.findRouter(tc.method, tc.path)
		if ok != tc.expectedOK {
			t.Errorf("Expected ok to be %v, got %v", tc.expectedOK, ok)
		}
		if node.node.path != tc.expectedPath {
			t.Errorf("Expected node path to be '%s', got '%s'", tc.expectedPath, node.node.path)
		}
	}
}
