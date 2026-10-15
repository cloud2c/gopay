package gopay

import "errors"

var (
	MissWechatInitParamErr = errors.New("missing wechat init parameter")
	MissAlipayInitParamErr = errors.New("missing alipay init parameter")
	MissPayPalInitParamErr = errors.New("missing paypal init parameter")
	MissAppleInitParamErr  = errors.New("missing apple init parameter")
	MissLakalaInitParamErr = errors.New("missing lakala init parameter")
	MissDouyinInitParamErr = errors.New("missing douyin init parameter")
	MissParamErr           = errors.New("missing required parameter")
	MarshalErr             = errors.New("marshal error")
	UnmarshalErr           = errors.New("unmarshal error")
	SignatureErr           = errors.New("signature error")
	VerifySignatureErr     = errors.New("verify signature error")
	// BusinessErr 表示渠道明确拒绝了业务请求（如参数不合法、余额不足）。
	// 各渠道的业务错误类型实现 errors.Is(err, gopay.BusinessErr)。
	BusinessErr     = errors.New("business error")
	CertNotMatchErr = errors.New("cert not match error")
	GetSignDataErr  = errors.New("get signature data error")
	BodyMapNilErr   = errors.New("body map is nil")
)
