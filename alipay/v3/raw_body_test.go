package alipay

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/cloud2c/gopay"
)

// 开了 AES 之后，验签要用的密文一度是存在 ClientV3 的字段上的，
// 而 ClientV3 是共享实例 —— 并发下会互相覆盖/清空，`-race` 直接报警。
// 这个用例锁住"密文按响应隔离"这件事，跑 -race 才有意义。
func TestRawBodyIsPerResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderTimestamp, "1")
		w.Header().Set(HeaderNonce, "n")
		w.Header().Set(HeaderSignature, "AA==")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ciphertext"))
	}))
	defer srv.Close()

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
	c.SetProxyHost(srv.URL)
	c.SetAESKey("okLL8PuZJQQzg8ldit/LGg==")
	c.aliPayPublicKey = &key.PublicKey // 让 autoVerifySignByCert 真的走进去

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, bs, err := c.doPost(context.Background(), gopay.BodyMap{"a": "b"}, "/v3/x", "")
			if err != nil {
				return
			}
			// 验签必然失败（签名是假的），这里只关心不出现数据竞争，
			// 以及每个响应拿到的都是自己那份密文
			_ = c.autoVerifySignByCert(res, bs)
			if got := string(rawBodyOf(res)); got != "ciphertext" {
				t.Errorf("密文串了或丢了: %q", got)
			}
		}()
	}
	wg.Wait()
}
