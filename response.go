package junyousdk

import (
	"net/http"
)

// Result SDK 结果结构（公共 API）
type Result[T any] struct {
	Code    int    `json:"code"`
	ErrCode string `json:"err_code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// newResult 创建结果（内部辅助函数）
func newResult[T any](code int, success bool, message string, data T) *Result[T] {
	return &Result[T]{
		Code:    code,
		Success: success,
		Message: message,
		Data:    data,
	}
}

// NewSuccessResult 创建成功结果
func NewSuccessResult[T any](message string, data T) *Result[T] {
	return newResult(http.StatusOK, true, message, data)
}

// NewSysErrorResult 创建系统错误结果
func NewSysErrorResult[T any](message string) *Result[T] {
	return newResult(http.StatusInternalServerError, false, message, *new(T))
}

// NewParamErrorResult 创建参数错误结果
func NewParamErrorResult[T any](message string) *Result[T] {
	return newResult(http.StatusBadRequest, false, message, *new(T))
}
