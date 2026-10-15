package wechat

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cloud2c/gopay"
)

// 测试用 RSA 私钥，只在内存中生成，不写入仓库。
func testPrivateKeyPEM(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	der := x509.MarshalPKCS1PrivateKey(key)
	return string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}))
}

// fakeTransport 返回预设的状态码与响应体，不访问网络。
type fakeTransport struct {
	status int
	body   string
	err    error
}

func (f fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &http.Response{
		StatusCode: f.status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(f.body)),
		Request:    req,
	}, nil
}

func newFakeWechatClient(t *testing.T, tr http.RoundTripper) *ClientV3 {
	t.Helper()
	c, err := NewClientV3("1900000109", "SERIAL0001", "0123456789abcdef0123456789abcdef", testPrivateKeyPEM(t))
	if err != nil {
		t.Fatal(err)
	}
	c.hc.SetTransport(tr)
	return c
}

func TestBusinessFailureIsTypedAndKeepsResponse(t *testing.T) {
	c := newFakeWechatClient(t, fakeTransport{status: http.StatusBadRequest, body: `{"code":"PARAM_ERROR","message":"参数错误"}`})

	rsp, err := c.V3ScheduledDeductPreSignApp(context.Background(), gopay.BodyMap{"appid": "wx0000000000000000"})
	if err == nil {
		t.Fatal("non-200 response must return an error")
	}
	if !errors.Is(err, gopay.BusinessErr) {
		t.Fatalf("error must satisfy errors.Is(err, gopay.BusinessErr): %v", err)
	}
	var biz *BizErr
	if !errors.As(err, &biz) || biz.Code != "PARAM_ERROR" || biz.HTTPStatus != http.StatusBadRequest {
		t.Fatalf("want *BizErr{PARAM_ERROR,400}, got %#v", err)
	}
	if rsp == nil || rsp.ErrResponse.Code != "PARAM_ERROR" {
		t.Fatalf("response must still be returned with ErrResponse, got %#v", rsp)
	}
}

func TestBusinessFailureWithUnparsableBodyStillTyped(t *testing.T) {
	c := newFakeWechatClient(t, fakeTransport{status: http.StatusInternalServerError, body: `<html>gateway</html>`})

	_, err := c.V3ScheduledDeductPreSignApp(context.Background(), gopay.BodyMap{"appid": "wx0000000000000000"})
	if !errors.Is(err, gopay.BusinessErr) {
		t.Fatalf("non-200 with unparsable body must still be a business error: %v", err)
	}
}

func TestTransportFailureIsNotBusinessError(t *testing.T) {
	c := newFakeWechatClient(t, fakeTransport{err: errors.New("connection refused")})

	_, err := c.V3ScheduledDeductPreSignApp(context.Background(), gopay.BodyMap{"appid": "wx0000000000000000"})
	if err == nil {
		t.Fatal("transport failure must return an error")
	}
	if errors.Is(err, gopay.BusinessErr) {
		t.Fatalf("transport failure must not be classified as a business error: %v", err)
	}
}

func TestBizErrMessageDoesNotEchoBody(t *testing.T) {
	err := bizErr(http.StatusForbidden, ErrResponse{Code: "NO_AUTH", Message: "无权限"})
	if strings.Contains(err.Error(), "{") {
		t.Fatalf("error text must not contain raw body markup: %q", err.Error())
	}
}
