package errs

import (
	"errors"
	"fmt"
)

// var (
//
//	ErrPointerOnly            = errors.New("orm: 只支持一级指针作为输入，例如 *User")
//	ErrNoRows                 = errors.New("orm: 未找到数据")
//	ErrTooManyReturnedColumns = errors.New("eorm: 过多列")
//
// )

// @RouterIsEmpty 10001
func RouterIsEmpty(col string) error {
	return fmt.Errorf("lychee-route: 路由是空字符串 router is emptry string : %s", col)
}

var (
	ErrRouterIsEmptry       = errors.New("web: 路由是空字符串")
	ErrRouterStartWithSlash = errors.New("web: 路由必须以 / 开头")
	ErrRouterEndWithSlash   = errors.New("web: 路由不能以 / 结尾")
	ErrRouterHadRoot        = errors.New("web: 路由冲突[/]")
	ErrRouterInvalid        = errors.New("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由")
	ErrRouterExisit         = errors.New("web: 路由冲突")
)

// func RouterInvalidPath(path string) error {
// 	return fmt.Errorf(, a ...any)
// }
