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
type matchNode struct {
	node  *node
	param map[string]string
}

func (m *matchNode) addValue(key string, value string) {
	if m.param == nil {
		m.param = map[string]string{key: value}
	}
	m.param[key] = value
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
		//root.route = path
		return
	}

	// /user/info  user/info
	segs := strings.Split(path[1:], "/")
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
	root.route = path

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

// findRouter retrun route node
func (r *router) findRouter(method string, path string) (*matchNode, bool) {
	node, ok := r.treeBranch[method]
	res := &matchNode{}
	if !ok {
		return nil, false
	}
	if path == "/" {
		res.node = node
		return res, true
	}
	segs := strings.Split(strings.Trim(path, "/"), "/")
	for _, s := range segs {
		var isParam bool
		node, isParam, ok = node.child(s)
		if !ok {
			return nil, false
		}
		if isParam {
			res.addValue(node.path[1:], s)
		}
	}
	res.node = node
	return res, true
}

// child
func (n *node) child(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
	}
	res, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
	}
	return res, false, true
}
