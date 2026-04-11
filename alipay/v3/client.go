package alipay

import (
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/go-pay/crypto/xpem"
	"github.com/go-pay/crypto/xrsa"
	"github.com/cloud2c/gopay"
	"github.com/cloud2c/gopay/pkg/xhttp"
	"github.com/go-pay/xlog"
)

// ClientV3 支付宝 V3
type ClientV3 struct {
	AppId              string
	AppCertSN          string
	AliPayPublicCertSN string
	AliPayRootCertSN   string
	AppAuthToken       string
	IsProd             bool
	aesKey             string // biz_content 加密的 AES KEY（Base64 编码）
	encryptType        string // 内容加密类型，默认 AES
	proxyHost          string // 代理host地址
	ivKey              []byte
	privateKey         *rsa.PrivateKey
	aliPayPublicKey    *rsa.PublicKey // 支付宝证书公钥内容 alipayPublicCert.crt
	DebugSwitch        gopay.DebugSwitch
	logger             xlog.XLogger
	requestIdFunc      xhttp.RequestIdHandler
	hc                 *xhttp.Client
	rawBodyForSign     []byte // 解密前的原始响应体，用于签名验证（支付宝 V3 签名是对密文做的）
}

// NewClientV3 初始化支付宝客户端 V3
// appid：应用ID
// privateKey：应用私钥，支持PKCS1和PKCS8
// isProd：是否是正式环境，沙箱环境请选择新版沙箱应用。
func NewClientV3(appid, privateKey string, isProd bool) (client *ClientV3, err error) {
	if appid == gopay.NULL || privateKey == gopay.NULL {
		return nil, gopay.MissAlipayInitParamErr
	}
	key := xrsa.FormatAlipayPrivateKey(privateKey)
	priKey, err := xpem.DecodePrivateKey([]byte(key))
	if err != nil {
		return nil, err
	}
	logger := xlog.NewLogger()
	logger.SetLevel(xlog.DebugLevel)
	client = &ClientV3{
		AppId:         appid,
		IsProd:        isProd,
		privateKey:    priKey,
		DebugSwitch:   gopay.DebugOff,
		logger:        logger,
		requestIdFunc: defaultRequestIdFunc,
		hc:            xhttp.NewClient(),
	}
	return client, nil
}

// 应用公钥证书内容设置 app_cert_sn、alipay_root_cert_sn、alipay_cert_sn
// appCertContent：应用公钥证书文件内容
// alipayRootCertContent：支付宝根证书文件内容
// alipayPublicCertContent：支付宝公钥证书文件内容
func (a *ClientV3) SetCert(appCertContent, alipayRootCertContent, alipayPublicCertContent []byte) (err error) {
	appCertSn, err := getCertSN(appCertContent)
	if err != nil {
		return fmt.Errorf("get app_cert_sn return err, but alse return alipay client. err: %w", err)
	}
	rootCertSn, err := getRootCertSN(alipayRootCertContent)
	if err != nil {
		return fmt.Errorf("get alipay_root_cert_sn return err, but alse return alipay client. err: %w", err)
	}
	publicCertSn, err := getCertSN(alipayPublicCertContent)
	if err != nil {
		return fmt.Errorf("get alipay_cert_sn return err, but alse return alipay client. err: %w", err)
	}

	// alipay public key
	pubKey, err := xpem.DecodePublicKey(alipayPublicCertContent)
	if err != nil {
		return fmt.Errorf("decode alipayPublicCertContent err: %w", err)
	}

	a.AppCertSN = appCertSn
	a.AliPayRootCertSN = rootCertSn
	a.AliPayPublicCertSN = publicCertSn
	a.aliPayPublicKey = pubKey
	return nil
}

// 设置自定义RequestId生成函数
func (a *ClientV3) SetRequestIdFunc(requestIdFunc xhttp.RequestIdHandler) *ClientV3 {
	if requestIdFunc != nil {
		a.requestIdFunc = requestIdFunc
	}
	return a
}

// 设置应用授权
func (a *ClientV3) SetAppAuthToken(appAuthToken string) *ClientV3 {
	a.AppAuthToken = appAuthToken
	return a
}

// SetBodySize 设置http response body size(MB)
func (a *ClientV3) SetBodySize(sizeMB int) *ClientV3 {
	if sizeMB > 0 {
		a.hc.SetBodySize(sizeMB)
	}
	return a
}

// SetHttpClient 设置自定义的xhttp.Client
func (a *ClientV3) SetHttpClient(client *xhttp.Client) *ClientV3 {
	if client != nil {
		a.hc = client
	}
	return a
}

// SetLogger 设置自定义的logger
func (a *ClientV3) SetLogger(logger xlog.XLogger) *ClientV3 {
	if logger != nil {
		a.logger = logger
	}
	return a
}

// SetAESKey 设置 V3 接口内容加密的 AES 密钥（Base64 编码）
// 设置此参数后，V3 POST 请求将自动对请求体进行 AES-128-CBC 加密，并添加 alipay-encrypt-type Header
// AES 密钥从支付宝开放平台「开发设置 > 接口内容加密方式」获取
func (a *ClientV3) SetAESKey(aesKey string) *ClientV3 {
	a.aesKey = aesKey
	a.encryptType = "AES"
	a.ivKey = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	return a
}

// SetEncryptType 设置内容加密类型，默认 AES
func (a *ClientV3) SetEncryptType(encryptType string) *ClientV3 {
	if encryptType != "" {
		a.encryptType = encryptType
	}
	return a
}

// IsEncryptEnabled 是否启用了内容加密
func (a *ClientV3) IsEncryptEnabled() bool {
	return a.aesKey != ""
}

// SetProxyHost 设置的 ProxyHost
// 使用场景：
// 1. 部署环境无法访问互联网，可以通过代理服务器访问
// 2. 不设置则默认 https://api.mch.weixin.qq.com
func (a *ClientV3) SetProxyHost(proxyHost string) *ClientV3 {
	before, found := strings.CutSuffix(proxyHost, "/")
	if found {
		a.proxyHost = before
		return a
	}
	a.proxyHost = proxyHost
	return a
}

// GetProxyHost 返回当前的 ProxyHost
func (a *ClientV3) GetProxyHost() string {
	return a.proxyHost
}
