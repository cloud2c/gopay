package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cloud2c/gopay"
	"github.com/go-pay/util"
	"github.com/go-pay/util/convert"
)

// v3 鉴权请求 Authorization Header
// bm: 请求参数（普通签名时使用，加密签名时传 nil）
// encryptedBody: 加密后的请求体字符串（加密签名时使用，普通签名时传空字符串）
// 当 encryptedBody 不为空时，使用密文字符串参与签名（先加密后签名）
func (a *ClientV3) authorization(method, uri string, bm gopay.BodyMap, appAuthToken, encryptedBody string) (string, error) {
	var (
		jb        = ""
		aat       = a.AppAuthToken // 默认值
		timestamp = convert.Int64ToString(time.Now().UnixNano() / int64(time.Millisecond))
		nonceStr  = util.RandomString(32)
		// app_id=2014060600164699,app_cert_sn=xxx,nonce=5f9fba93-bbb2-40f0-b328-04d5ead3e131,timestamp=1667804301218
		authString = "app_id=" + a.AppId + ",app_cert_sn=" + a.AppCertSN + ",nonce=" + nonceStr + ",timestamp=" + timestamp
	)
	if a.AppCertSN == gopay.NULL {
		authString = "app_id=" + a.AppId + ",nonce=" + nonceStr + ",timestamp=" + timestamp
	}
	if appAuthToken != "" {
		aat = appAuthToken
	}

	// 确定签名用的 body 内容
	if encryptedBody != "" {
		// 加密请求：body 直接使用密文字符串
		jb = encryptedBody
	} else if bm != nil {
		// 普通请求：签名body里需要删掉 alipay-app-auth-token
		bm.Remove(HeaderAppAuthToken)
		jb = bm.JsonBody()
	}

	// ${authString}\n	步骤1中生成的认证串 authString。
	// ${httpMethod}\n	本次请求的 http 方法，例如 GET\POST\PUT 等。
	// ${httpReuqestUrl}\n   本次请求的 uri 信息，包括 queryString，不包括域名，例如 /v3/alipay/marketing/activity/ordervoucher/get?id=123。
	// ${httpRequestBody}\n	本次请求的 body 内容。当使用GET等请求时，body 为空，该值传入空字符串，即""。
	// ${appAuthToken}\n		应用授权令牌，和 header 参数中 alipay-app-auth-token 值保持一致。可选参数，不使用代调用模式时，不需要传入该字段。
	//
	// 示例：
	// app_id=2014060600164699,timestamp=1655869956477,nonce=eb4ade8f-8cfa-4ebf-a048-7eb52684ab32,expired_seconds=120
	// POST
	// /v3/alipay/marketing/activity/ordervoucher/create?auth_token=123
	// {"activity_name": "单品特价满10减1活动","publish_start_time": "2022-02-01 00:00:01"}
	//
	// body 空示例：
	// app_id=2014060600164699,timestamp=1655869956477,nonce=eb4ade8f-8cfa-4ebf-a048-7eb52684ab32,expired_seconds=120
	// GET
	// /v3/alipay/marketing/activity/ordervoucher?id=123
	//
	// 代调示例：
	// app_id=2014060600164699,timestamp=1655869956477,nonce=eb4ade8f-8cfa-4ebf-a048-7eb52684ab32,expired_seconds=120
	// GET
	// /v3/alipay/marketing/activity/ordervoucher?id=123
	// 202212BB_D64b2be8afd4b6c8468cf585bd05E50
	signStr := authString + "\n" + method + "\n" + uri + "\n" + jb + "\n"
	if aat != gopay.NULL {
		signStr += aat + "\n"
	}
	if a.DebugSwitch == gopay.DebugOn {
		if encryptedBody != "" {
			a.logger.Debugf("Alipay_V3_Encrypt_SignString:\n%s", signStr)
		} else {
			a.logger.Debugf("Alipay_V3_SignString:\n%s", signStr)
		}
	}

	sign, err := a.rsaSign(signStr)
	if err != nil {
		return "", err
	}
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Sign:\n%s", sign)
	}
	// authorization: ${签名算法} ${authString},sign=${signature}
	authorization := SignTypeRSA + " " + authString + ",sign=" + sign
	return authorization, nil
}

func (a *ClientV3) rsaSign(str string) (string, error) {
	if a.privateKey == nil {
		return "", errors.New("privateKey can't be nil")
	}
	h := sha256.New()
	h.Write([]byte(str))
	result, err := rsa.SignPKCS1v15(rand.Reader, a.privateKey, crypto.SHA256, h.Sum(nil))
	if err != nil {
		return gopay.NULL, fmt.Errorf("[%w]: %+v", gopay.SignatureErr, err)
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

// =============================== 同步验签 ===============================

func (a *ClientV3) autoVerifySignByCert(res *http.Response, body []byte) (err error) {
	if a.aliPayPublicKey == nil {
		// 没有支付宝公钥 = 没法验签。默认沿用历史行为直接放过，
		// 但这意味着响应是**未经验证**的，调用方可以用 SetVerifyRequired(true) 要求拦下来。
		if a.verifyRequired {
			return fmt.Errorf("[%w]: alipay public key not set, call SetCert (cert mode) "+
				"or SetAliPayPublicKey (key mode) first; "+
				"pass SetVerifyRequired(false) only if you accept unverified responses",
				gopay.VerifySignatureErr)
		}
		// 走到这里说明调用方显式关掉了验签要求，响应是**未经验证**的
		return nil
	}
	ts := res.Header.Get(HeaderTimestamp)
	nonce := res.Header.Get(HeaderNonce)
	sign := res.Header.Get(HeaderSignature)

	// 证书模式下支付宝会回 alipay-sn，标明它本次用的是哪本证书。
	// 官方要求商家确保本地支付宝证书与它一致，不一致需要更新支付宝公钥证书。
	//
	// 单独报出来，而不是让它掉进下面的验签失败：证书轮换后本地公钥必然对不上，
	// 报「签名不匹配」会把人引去查密钥和报文，而实际要做的是换一本证书。
	//
	// 两边任一为空就跳过：公钥（密钥）模式下支付宝不回这个头，
	// AliPayPublicCertSN 也为空，不该拿它当失败条件。
	if aliSN := res.Header.Get(HeaderAliPaySN); aliSN != gopay.NULL &&
		a.AliPayPublicCertSN != gopay.NULL && aliSN != a.AliPayPublicCertSN {
		return fmt.Errorf("[%w]: alipay cert sn mismatch, response=%s, local=%s. "+
			"请更新本地支付宝公钥证书", gopay.VerifySignatureErr, aliSN, a.AliPayPublicCertSN)
	}
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_VerifySignHeader: alipay-timestamp=[%s], alipay-nonce=[%s], alipay-signature=[%s]", ts, nonce, sign)
	}
	// 支付宝 V3 签名是对加密后的密文做的，当设置了 AES Key 时，
	// body 已被 doPost 解密成明文，验签要用挂在响应上的那份密文
	signBody := body
	if raw := rawBodyOf(res); a.aesKey != "" && raw != nil {
		signBody = raw
	}
	signData := ts + "\n" + nonce + "\n" + string(signBody) + "\n"

	// 不吞这个错误：sign 头缺失或被截断时，解出来是空/半截 bytes，
	// 后面 Verify 一样会失败，但报的是"签名不匹配"，
	// 让人去查密钥和报文，而真正的原因是压根没收到签名
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return fmt.Errorf("[%w]: decode alipay-signature header err: %v", gopay.VerifySignatureErr, err)
	}
	sum256 := sha256.Sum256([]byte(signData))
	if err = rsa.VerifyPKCS1v15(a.aliPayPublicKey, crypto.SHA256, sum256[:], signBytes); err != nil {
		return fmt.Errorf("[%w]: %v", gopay.VerifySignatureErr, err)
	}
	return nil
}
