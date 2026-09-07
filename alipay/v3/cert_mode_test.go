package alipay

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloud2c/gopay"
)

func newTestClientV3(t *testing.T, srvURL string) (*ClientV3, *rsa.PrivateKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	// NewClientV3 内部会补 PEM 头，这里给裸 base64
	rawKey := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(key))
	c, err := NewClientV3("2021000000000000", rawKey, true)
	if err != nil {
		t.Fatal(err)
	}
	c.SetProxyHost(srvURL)
	return c, key
}

// 官方「签名规则」要求证书模式的请求头带 alipay-root-cert-sn：
//
//	Authorization: ${签名算法} ${authString},sign=${signature}
//	alipay-root-cert-sn: ${alipayRootCertSn}
//
// 之前只发了 Authorization，根证书序列号一个字节都没发。
func TestRootCertSNHeaderSentInCertMode(t *testing.T) {
	var gotSN string
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSN = r.Header.Get(HeaderRootCertSN)
		gotAuth = r.Header.Get(HeaderAuthorization)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
	}))
	defer srv.Close()

	c, _ := newTestClientV3(t, srv.URL)
	c.AliPayRootCertSN = "root-sn-123"
	c.AppCertSN = "app-sn-456"

	if _, _, err := c.doPost(context.Background(), gopay.BodyMap{"a": "b"}, "/v3/x", ""); err != nil {
		t.Fatal(err)
	}
	if gotSN != "root-sn-123" {
		t.Errorf("alipay-root-cert-sn 没发出去: got %q", gotSN)
	}
	// app_cert_sn 走 authString，不是独立的头
	if want := "app_cert_sn=app-sn-456"; !contains(gotAuth, want) {
		t.Errorf("authString 里没有 %s: %q", want, gotAuth)
	}
}

// 密钥（公钥）模式没有根证书，不该发这个头
func TestRootCertSNHeaderAbsentInKeyMode(t *testing.T) {
	var present bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, present = r.Header[http.CanonicalHeaderKey(HeaderRootCertSN)]
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
	}))
	defer srv.Close()

	c, _ := newTestClientV3(t, srv.URL)
	if _, _, err := c.doPost(context.Background(), gopay.BodyMap{"a": "b"}, "/v3/x", ""); err != nil {
		t.Fatal(err)
	}
	if present {
		t.Error("密钥模式不该发 alipay-root-cert-sn")
	}
}

// 支付宝轮换证书后 alipay-sn 会和本地对不上。
// 这时要直接说「换证书」，而不是掉进验签失败报「签名不匹配」——
// 后者会把排查引向密钥和报文。
func TestAliPaySNMismatchReportsCertError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderTimestamp, "1")
		w.Header().Set(HeaderNonce, "n")
		w.Header().Set(HeaderSignature, "AA==")
		w.Header().Set(HeaderAliPaySN, "new-cert-sn")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
	}))
	defer srv.Close()

	c, key := newTestClientV3(t, srv.URL)
	c.aliPayPublicKey = &key.PublicKey
	c.AliPayPublicCertSN = "old-cert-sn"

	res, bs, err := c.doPost(context.Background(), gopay.BodyMap{"a": "b"}, "/v3/x", "")
	if err != nil {
		t.Fatal(err)
	}
	err = c.autoVerifySignByCert(res, bs)
	if err == nil {
		t.Fatal("证书号不一致却放过了")
	}
	if !contains(err.Error(), "cert sn mismatch") {
		t.Errorf("错误信息没点明证书号不一致: %v", err)
	}
}

// SetAliPayPublicKey 是密钥模式的验签入口。
// 在此之前 v3 只有 SetCert 会设 aliPayPublicKey，密钥模式下验签被静默跳过。
func TestSetAliPayPublicKeyEnablesVerification(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	rawKey := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(key))
	c, err := NewClientV3("2021000000000000", rawKey, true)
	if err != nil {
		t.Fatal(err)
	}
	if c.aliPayPublicKey != nil {
		t.Fatal("初始不该有支付宝公钥")
	}
	// 支付宝公钥是 PKIX(X.509 SubjectPublicKeyInfo) 编码，不是 PKCS1
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pub := base64.StdEncoding.EncodeToString(der)
	if err := c.SetAliPayPublicKey([]byte(pub)); err != nil {
		t.Fatalf("裸 base64 公钥应当被接受: %v", err)
	}
	if c.aliPayPublicKey == nil {
		t.Fatal("SetAliPayPublicKey 没设上")
	}
	if err := c.SetAliPayPublicKey(nil); err == nil {
		t.Error("空公钥应当报错")
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
