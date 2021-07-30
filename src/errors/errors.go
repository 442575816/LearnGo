package errors

import "errors"

var (
	// HandleAlreadyExists 处理器已经存在
	HandleAlreadyExists = errors.New("handler exists already")
	// NotSupport 不支持该操作
	NotSupport = errors.New("not support this operation")
)
