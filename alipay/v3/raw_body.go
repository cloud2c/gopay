package alipay

import (
	"context"
	"net/http"
)

// 开了内容加密（SetAESKey）之后，支付宝 V3 的响应签名是对**密文**做的，
// 而 doPost 会就地把 bs 解密成明文再返回，验签时拿不到密文了。
//
// 这里把密文挂在这次响应自己的 Request 上，而不是存回 ClientV3。
//
// ⚠️ 早先的写法是 client 上一个 rawBodyForSign 字段：doPost 写、
// autoVerifySignByCert 读完清空。ClientV3 是长期共享实例（典型用法就是放在
// 服务的依赖容器里给所有请求用），这个字段全程无锁，`go test -race` 直接报
// 两处 DATA RACE。而且不只是竞态告警，业务结果也是错的：
//
//	并发请求 A、B 同时回来
//	B 的密文覆盖掉 A 的  ->  A 拿 B 的 body 验签  ->  验签失败
//	A 读完把字段清空     ->  B 退回用明文验签      ->  验签失败
//
// 两边都是「签名验证失败」，而报文本身完全正常 —— 这种错在生产上极难定位，
// 量小时还偶尔能过，一上并发就开始随机失败。
//
// 挂在 res.Request 上则天然按响应隔离：http.Client.Do 保证返回的 Response.Request
// 非空，且它随响应一起被 GC，不需要清理，也不会串到下一次请求。
type rawBodyCtxKey struct{}

// stashRawBody 记下解密前的密文，供 autoVerifySignByCert 验签使用
func stashRawBody(res *http.Response, raw []byte) {
	if res == nil || res.Request == nil || len(raw) == 0 {
		return
	}
	buf := make([]byte, len(raw))
	copy(buf, raw)
	res.Request = res.Request.WithContext(context.WithValue(res.Request.Context(), rawBodyCtxKey{}, buf))
}

// rawBodyOf 取回密文，没有则返回 nil（未开加密、或响应未走加密分支）
func rawBodyOf(res *http.Response) []byte {
	if res == nil || res.Request == nil {
		return nil
	}
	raw, _ := res.Request.Context().Value(rawBodyCtxKey{}).([]byte)
	return raw
}
