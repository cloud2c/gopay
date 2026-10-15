package wechat

import (
	"fmt"

	"github.com/cloud2c/gopay"
)

// BizErr 表示微信 V3 接口返回了业务失败（HTTP 状态码非 200）。
// 响应结构体中的 ErrResponse 字段仍然保留，调用方可以继续读取；
// 错误本身只包含状态码、错误码与错误描述，不包含响应体原文。
//
// 与支付宝的 BizErr 一样，满足 errors.Is(err, gopay.BusinessErr)，
// 使只检查 err 的调用方不会把失败当成成功。
type BizErr struct {
	HTTPStatus int
	Code       string
	Message    string
}

func (e *BizErr) Error() string {
	return fmt.Sprintf("wechat v3: http %d, code %s: %s", e.HTTPStatus, e.Code, e.Message)
}

// Is 使微信业务错误满足 errors.Is(err, gopay.BusinessErr)。
func (e *BizErr) Is(target error) bool {
	return target == gopay.BusinessErr
}

// bizErr 由非 200 响应构造 BizErr。
func bizErr(status int, e ErrResponse) error {
	return &BizErr{HTTPStatus: status, Code: e.Code, Message: e.Message}
}
