package lychee

import (
	"lychee/interal/errs"
	"strings"
)

// [path]->node
type router struct {
	treeBranch map[string]*node
}

type node struct {
	path        string
	route       string
	children    map[string]*node
	paramChild  *node
	contextFunc ContextFunc //context
}

func NewRouter() router {
	return router{
		treeBranch: map[string]*node{},
	}
}
func (r *router) addRoute(method string, path string, cfunc ContextFunc) {
	path = checkRoute(path)
	root, ok := r.treeBranch[method]
	if !ok {
		root = &node{path: "/"}
		r.treeBranch[method] = root //add
	}
	if path == "/" {
		if root.contextFunc != nil {
			panic(errs.ErrRouterHadRoot)
		}
		root.contextFunc = cfunc
		return
	}

	str := path[1:] // /user/info  user/info
	segs := strings.Split(str, "/")
	for _, seg := range segs {
		if seg == "" {
			panic(errs.ErrRouterInvalid)
		}
		root = root.childCheck(seg)
	}
	if root.contextFunc != nil {
		panic(errs.ErrRouterExisit)
	}
	root.contextFunc = cfunc
	root.path = path

}

func (n *node) childCheck(path string) *node {
	if path[0] == ':' {
		if n.paramChild != nil {
			if n.paramChild.path != path {
				panic(errs.ErrRouterExisit)
			}
		} else {
			n.paramChild = &node{path: path}
		}
		return n.paramChild
	}
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}

func checkRoute(path string) string {
	if path == "" {
		panic(errs.RouterIsEmpty(path))
	}
	if path[0] != '/' {
		panic(errs.ErrRouterStartWithSlash)
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic(errs.ErrRouterEndWithSlash)
	}
	return path
}
