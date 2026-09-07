package alipay

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"

	"github.com/cloud2c/gopay"
	"github.com/cloud2c/gopay/pkg/xhttp"
	"github.com/go-pay/crypto/xpem"
	"github.com/go-pay/crypto/xrsa"
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
	verifyRequired     bool // 见 SetVerifyRequired
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
//
// 注意：不同 API 对加密的要求不同：
//   - 人脸核身/OCR/数据核验等 V2 证书模式接口：必须使用 AES 加密 biz_content
//   - 扫码支付/交易创建等 V3 接口：不支持内容加密，设置了反而会导致参数错误
//
// 使用建议：
// 方案 A（推荐）：创建两个独立的 Client 实例
//
//	faceClient := alipay.NewClientV3(...).SetAESKey(aesKey)  // 人脸核身专用
//	payClient := alipay.NewClientV3(...)                      // 支付专用（不设置 AES）
//
// 方案 B：使用 WithoutAES() 临时获取一个不加密的 Client（线程安全）
//
//	client.SetAESKey(aesKey)
//	client.FaceVerificationInitialize(...)       // 使用 AES 加密
//	client.WithoutAES().TradeCreate(...)          // 获取无加密的 Client 副本调用支付
//	client.FaceVerificationQuery(...)             // 原 Client 不受影响，仍启用 AES
func (a *ClientV3) SetAESKey(aesKey string) *ClientV3 {
	a.aesKey = aesKey
	a.encryptType = "AES"
	a.ivKey = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	return a
}

// SetAliPayPublicKey 公钥模式：直接设置**支付宝公钥**（不是证书），用于同步响应验签。
//
// 与 SetCert 是二选一的两种模式：
//
//	证书模式  SetCert()             app_cert_sn 进签名串，验签用证书里的公钥
//	公钥模式  SetAliPayPublicKey()  不带 app_cert_sn，验签用这里给的公钥
//
// 在此之前 v3 只有 SetCert 会给 aliPayPublicKey 赋值，公钥模式下
// autoVerifySignByCert 会直接 return nil —— 也就是**响应完全不验签**，
// 而且是静默的。补上这个入口后，公钥模式才真正能验签。
//
// publicKey 支持带 PEM 头的完整公钥，也支持支付宝控制台直接拷出来的裸 base64
// （内部会补 PEM 头）。
//
// 注意别和 SetCert 同时调：后调用的会覆盖 aliPayPublicKey，而 AppCertSN
// 只有 SetCert 会设，混着用会得到「带 app_cert_sn 但用普通公钥验签」这种
// 支付宝不认的组合。
func (a *ClientV3) SetAliPayPublicKey(publicKey []byte) (err error) {
	if len(publicKey) == 0 {
		return errors.New("alipay public key is empty")
	}
	key := publicKey
	// 裸 base64 补上 PEM 头，和 NewClientV3 处理私钥的方式一致
	if !bytes.Contains(key, []byte("-----BEGIN")) {
		key = []byte(xrsa.FormatAlipayPublicKey(string(key)))
	}
	pubKey, err := xpem.DecodePublicKey(key)
	if err != nil {
		return fmt.Errorf("decode alipay public key err: %w", err)
	}
	a.aliPayPublicKey = pubKey
	return nil
}

// SetVerifyRequired 要求同步响应必须验签，没有支付宝公钥就直接报错。
//
// 默认是 false，也就是**未配置证书时 autoVerifySignByCert 什么都不做、返回 nil**。
// 那是历史行为，为了不打断只用公钥模式的调用方，这里没有直接改默认值。
//
// 但要清楚这个默认值的含义：v3 里给 aliPayPublicKey 赋值的入口**只有 SetCert**，
// 没有别的方法能设支付宝公钥。所以「公钥模式」在 v3 上等价于**响应完全不验签**，
// 而且是静默的 —— 调用方拿到的 err 是 nil，看起来一切正常。
//
// 对接支付这种场景建议显式打开：
//
//	client, _ := alipay.NewClientV3(appId, privateKey, isProd)
//	client.SetCert(appCert, rootCert, alipayCert)
//	client.SetVerifyRequired(true)   // 漏配证书时立刻暴露，而不是一路不验签
//
// 必须在并发使用之前设置好（和 SetCert 一样，属于初始化期配置）。
func (a *ClientV3) SetVerifyRequired(required bool) *ClientV3 {
	a.verifyRequired = required
	return a
}

// WithoutAES 返回一个新的 Client 副本，不启用 AES 加密
// 原 Client 实例不受任何影响，此方法线程安全
//
// 使用场景：当全局设置了 AES Key（用于人脸核身等接口），但需要调用不支持加密的支付接口时
//
// 示例：
//
//	client.SetAESKey(aesKey)
//	client.FaceVerificationInitialize(...)       // 使用 AES 加密
//	client.WithoutAES().TradeCreate(...)          // 使用无加密的副本调用支付
//	client.FaceVerificationQuery(...)             // 原 Client 仍启用 AES
//
// 注意：此方法返回的是新实例，不会修改原 Client 的状态
func (a *ClientV3) WithoutAES() *ClientV3 {
	return a.Clone()
}

// Clone 克隆当前 Client 实例，返回一个独立的新实例
// 新实例共享私钥和证书配置，但不继承 AES 加密配置
//
// 使用场景：需要同时支持加密接口（人脸核身）和非加密接口（支付）时
//
// 示例：
//
//	baseClient := alipay.NewClientV3(appid, privateKey, isProd).SetCert(...)
//
//	// 人脸核身专用 client（带 AES 加密）
//	faceClient := baseClient.Clone().SetAESKey(aesKey)
//	faceClient.FaceVerificationInitialize(...)
//
//	// 支付专用 client（不带 AES 加密）
//	payClient := baseClient.Clone()
//	payClient.TradeCreate(...)
func (a *ClientV3) Clone() *ClientV3 {
	return &ClientV3{
		AppId:              a.AppId,
		AppCertSN:          a.AppCertSN,
		AliPayPublicCertSN: a.AliPayPublicCertSN,
		AliPayRootCertSN:   a.AliPayRootCertSN,
		AppAuthToken:       a.AppAuthToken,
		IsProd:             a.IsProd,
		aesKey:             "", // 克隆时不继承 AES Key，避免误用
		encryptType:        "",
		proxyHost:          a.proxyHost,
		verifyRequired:     a.verifyRequired,
		privateKey:         a.privateKey,
		aliPayPublicKey:    a.aliPayPublicKey,
		DebugSwitch:        a.DebugSwitch,
		logger:             a.logger,
		requestIdFunc:      a.requestIdFunc,
		hc:                 a.hc, // 共享 HTTP Client
	}
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
// 2. 不设置则默认 https://openapi.alipay.com
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
