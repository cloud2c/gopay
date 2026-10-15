package xhttp

import (
	"net/http"
	"testing"
)

// 默认客户端必须校验服务端证书，防止支付请求在未校验的 TLS 连接上传输。
func TestDefaultClientVerifiesServerCertificate(t *testing.T) {
	c := NewClient()
	tr, ok := c.HttpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("default transport is %T, want *http.Transport", c.HttpClient.Transport)
	}
	if tr.TLSClientConfig == nil {
		t.Fatal("TLSClientConfig is nil; default verification must be explicit")
	}
	if tr.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("default client skips server certificate verification")
	}
}

// SetTLSVerify(false) 只能显式关闭校验，且可再次开启。
func TestSetTLSVerifyTogglesExplicitly(t *testing.T) {
	c := NewClient().SetTLSVerify(false)
	tr := c.HttpClient.Transport.(*http.Transport)
	if !tr.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("SetTLSVerify(false) did not disable verification")
	}
	c.SetTLSVerify(true)
	if tr.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("SetTLSVerify(true) did not re-enable verification")
	}
}
