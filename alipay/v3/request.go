package alipay

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-pay/crypto/aes"
	"github.com/cloud2c/gopay"
	"github.com/cloud2c/gopay/pkg/xhttp"
	"github.com/go-pay/util"
)

var defaultRequestIdFunc = &requestIdFunc{}

type requestIdFunc struct{}

func (d *requestIdFunc) RequestId() string {
	return fmt.Sprintf("%s-%d", util.RandomString(21), time.Now().Unix())
}

// DoAliPayAPISelfV3 支付宝接口自行实现方法
func (a *ClientV3) DoAliPayAPISelfV3(ctx context.Context, method, path string, bm gopay.BodyMap, aliRsp any) (res *http.Response, err error) {
	var (
		bs []byte
	)
	aat := bm.GetString(HeaderAppAuthToken)
	switch method {
	case MethodGet:
		bm.Remove(HeaderAppAuthToken)
		uri := path + "?" + bm.EncodeURLParams()
		res, bs, err = a.doGet(ctx, uri, aat)
		if err != nil {
			return nil, err
		}
	case MethodPost:
		res, bs, err = a.doPost(ctx, bm, path, aat)
		if err != nil {
			return nil, err
		}
	case MethodPatch:
		res, bs, err = a.doPatch(ctx, bm, path, aat)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("method:%s not support", method)
	}
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	return res, nil
}

// doPost V3 POST 请求，自动判断是否需要加密
// 当设置了 AES Key 时，自动执行：加密请求体 → 基于密文签名（先加密后签名）→ 发送加密请求
// 未设置 AES Key 时，执行普通签名流程
func (a *ClientV3) doPost(ctx context.Context, bm gopay.BodyMap, uri, aat string) (res *http.Response, bs []byte, err error) {
	var url = v3BaseUrlCh + uri
	if !a.IsProd {
		url = v3SandboxBaseUrl + uri
	}
	if a.proxyHost != "" {
		url = a.proxyHost + uri
	}

	var (
		sendBody      string
		authorization string
	)

	// V3 内容加密：当设置了 AES Key 时，先加密后签名
	if a.aesKey != "" && bm != nil {
		// 步骤 1：加密请求体
		encryptedBody, encErr := a.encryptV3Body(bm.JsonBody())
		if encErr != nil {
			return nil, nil, fmt.Errorf("encrypt request body error: %w", encErr)
		}
		sendBody = encryptedBody
		// 步骤 2：基于密文字符串计算签名（先加密后签名）
		authorization, err = a.authorization(MethodPost, uri, nil, aat, encryptedBody)
		if err != nil {
			return nil, nil, fmt.Errorf("authorization with encrypted body error: %w", err)
		}
	} else {
		// 普通签名流程
		if bm != nil {
			sendBody = bm.JsonBody()
		}
		authorization, err = a.authorization(MethodPost, uri, bm, aat, "")
		if err != nil {
			return nil, nil, err
		}
	}

	req := a.hc.Req() // default json
	req.Header.Add(HeaderAuthorization, authorization)
	req.Header.Add(HeaderRequestID, a.requestIdFunc.RequestId())
	req.Header.Add(HeaderSdkVersion, "gopay/"+gopay.Version)
	// V3 内容加密：添加 alipay-encrypt-type Header，并将 Content-Type 设为 text/plain
	if a.aesKey != "" {
		req.Header.Add(HeaderEncryptType, a.encryptType)
		req.Header.Set("Content-Type", "text/plain")
	}
	if aat != gopay.NULL {
		req.Header.Add(HeaderAppAuthToken, aat)
	} else if a.AppAuthToken != "" {
		req.Header.Add(HeaderAppAuthToken, a.AppAuthToken)
	}
	req.Header.Add("Accept", "application/json")
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Url: %s", url)
		if a.aesKey != "" && bm != nil {
			a.logger.Debugf("Alipay_V3_Origin_Body: %s", bm.JsonBody())
			a.logger.Debugf("Alipay_V3_Encrypt_Body: %s", sendBody)
		} else {
			a.logger.Debugf("Alipay_V3_Req_Body: %s", sendBody)
		}
		a.logger.Debugf("Alipay_V3_Req_Headers: %#v", req.Header)
	}
	res, bs, err = req.Post(url).SendString(sendBody).EndBytesForAlipayV3(ctx)
	if err != nil {
		return nil, nil, err
	}

	// V3 内容解密：如果响应包含加密标识，则解密响应体
	if res.StatusCode == http.StatusOK && a.aesKey != "" {
		if res.Header.Get(HeaderContentEncrypt) != "" {
			decryptedBs, decErr := a.decryptV3Body(string(bs))
			if decErr != nil {
				a.logger.Debugf("Alipay_V3_Decrypt_Response_Error: %v", decErr)
			} else {
				bs = []byte(decryptedBs)
			}
		}
	}

	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Response: %d >> %s", res.StatusCode, string(bs))
		a.logger.Debugf("Alipay_V3_Rsp_Headers: %#v", res.Header)
	}
	return res, bs, nil
}

func (a *ClientV3) doPatch(ctx context.Context, bm gopay.BodyMap, uri, aat string) (res *http.Response, bs []byte, err error) {
	authorization, err := a.authorization(MethodPatch, uri, bm, aat, "")
	if err != nil {
		return nil, nil, err
	}

	var url = v3BaseUrlCh + uri
	if !a.IsProd {
		url = v3SandboxBaseUrl + uri
	}
	if a.proxyHost != "" {
		url = a.proxyHost + uri
	}
	req := a.hc.Req() // default json
	req.Header.Add(HeaderAuthorization, authorization)
	req.Header.Add(HeaderRequestID, a.requestIdFunc.RequestId())
	req.Header.Add(HeaderSdkVersion, "gopay/"+gopay.Version)
	if aat != gopay.NULL {
		req.Header.Add(HeaderAppAuthToken, aat)
	} else if a.AppAuthToken != "" {
		req.Header.Add(HeaderAppAuthToken, a.AppAuthToken)
	}
	req.Header.Add("Accept", "application/json")
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Url: %s", url)
		a.logger.Debugf("Alipay_V3_Req_Body: %s", bm.JsonBody())
		a.logger.Debugf("Alipay_V3_Req_Headers: %#v", req.Header)
	}
	res, bs, err = req.Patch(url).SendBodyMap(bm).EndBytesForAlipayV3(ctx)
	if err != nil {
		return nil, nil, err
	}

	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Response: %d >> %s", res.StatusCode, string(bs))
		a.logger.Debugf("Alipay_V3_Rsp_Headers: %#v", res.Header)
	}
	return res, bs, nil
}

func (a *ClientV3) doPut(ctx context.Context, bm gopay.BodyMap, uri, aat string) (res *http.Response, bs []byte, err error) {
	authorization, err := a.authorization(MethodPut, uri, bm, aat, "")
	if err != nil {
		return nil, nil, err
	}

	var url = v3BaseUrlCh + uri
	if !a.IsProd {
		url = v3SandboxBaseUrl + uri
	}
	if a.proxyHost != "" {
		url = a.proxyHost + uri
	}
	req := a.hc.Req() // default json
	req.Header.Add(HeaderAuthorization, authorization)
	req.Header.Add(HeaderRequestID, a.requestIdFunc.RequestId())
	req.Header.Add(HeaderSdkVersion, "gopay/"+gopay.Version)
	if aat != gopay.NULL {
		req.Header.Add(HeaderAppAuthToken, aat)
	} else if a.AppAuthToken != "" {
		req.Header.Add(HeaderAppAuthToken, a.AppAuthToken)
	}
	req.Header.Add("Accept", "application/json")
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Url: %s", url)
		a.logger.Debugf("Alipay_V3_Req_Body: %s", bm.JsonBody())
		a.logger.Debugf("Alipay_V3_Req_Headers: %#v", req.Header)
	}
	res, bs, err = req.Put(url).SendBodyMap(bm).EndBytesForAlipayV3(ctx)
	if err != nil {
		return nil, nil, err
	}

	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Response: %d > %s", res.StatusCode, string(bs))
		a.logger.Debugf("Alipay_V3_Rsp_Headers: %#v", res.Header)
	}
	return res, bs, nil
}

func (a *ClientV3) doGet(ctx context.Context, uri, aat string) (res *http.Response, bs []byte, err error) {
	authorization, err := a.authorization(MethodGet, uri, nil, aat, "")
	if err != nil {
		return nil, nil, err
	}

	var url = v3BaseUrlCh + uri
	if !a.IsProd {
		url = v3SandboxBaseUrl + uri
	}
	if a.proxyHost != "" {
		url = a.proxyHost + uri
	}
	req := a.hc.Req() // default json
	req.Header.Add(HeaderAuthorization, authorization)
	req.Header.Add(HeaderRequestID, a.requestIdFunc.RequestId())
	req.Header.Add(HeaderSdkVersion, "gopay/"+gopay.Version)
	if aat != gopay.NULL {
		req.Header.Add(HeaderAppAuthToken, aat)
	} else if a.AppAuthToken != "" {
		req.Header.Add(HeaderAppAuthToken, a.AppAuthToken)
	}
	req.Header.Add("Accept", "application/json")
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Url: %s", url)
		a.logger.Debugf("Alipay_V3_Req_Headers: %#v", req.Header)
	}
	res, bs, err = req.Get(url).EndBytesForAlipayV3(ctx)
	if err != nil {
		return nil, nil, err
	}

	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Response: %d > %s", res.StatusCode, string(bs))
		a.logger.Debugf("Alipay_V3_Rsp_Headers: %#v", res.Header)
	}
	return res, bs, nil
}

// doProdPostFile 文件上传 POST 请求
// 自动处理文件上传的签名逻辑：从 bm 中分离 file 字段 → 非文件参数编码到 data 字段 → 基于 data 签名 → 将 file 加回 bm → 发送 multipart 请求
// 调用方只需传入包含所有参数（含 file）的 BodyMap，无需手动处理 data 编码和签名
func (a *ClientV3) doProdPostFile(ctx context.Context, bm gopay.BodyMap, uri, aat string) (res *http.Response, bs []byte, err error) {
	// 步骤 1：提取 file 字段到临时 map，并收集非 file 参数用于签名
	tempFile := make(gopay.BodyMap)
	signMap := make(gopay.BodyMap)
	bm.Range(func(k string, v any) bool {
		if file, ok := v.(*gopay.File); ok {
			tempFile.SetFormFile(k, file)
			bm.Remove(k)
		} else {
			signMap.Set(k, v)
		}
		return true
	})

	// 步骤 2：将非 file 参数编码到 data 字段中
	bm.SetBodyMap("data", func(b gopay.BodyMap) {
		signMap.Range(func(k string, v any) bool {
			b.Set(k, v)
			return true
		})
	})
	// 从 bm 顶层移除非 file 参数（已在 data 中）
	signMap.Range(func(k string, v any) bool {
		bm.Remove(k)
		return true
	})

	// 步骤 3：使用 signMap 计算签名（签名时排除 file 和 data 字段，仅用原始非 file 参数）
	authorization, err := a.authorization(MethodPost, uri, signMap, aat, "")
	if err != nil {
		return nil, nil, err
	}

	// 步骤 4：将 file 字段加回 bm（用于 multipart 发送）
	tempFile.Range(func(k string, v any) bool {
		bm.SetFormFile(k, v.(*gopay.File))
		return true
	})

	var url = v3BaseUrlCh + uri
	if a.proxyHost != "" {
		url = a.proxyHost + uri
	}
	req := a.hc.Req(xhttp.TypeMultipartFormData)
	req.Header.Add(HeaderAuthorization, authorization)
	req.Header.Add(HeaderRequestID, a.requestIdFunc.RequestId())
	req.Header.Add(HeaderSdkVersion, "gopay/"+gopay.Version)
	if aat != gopay.NULL {
		req.Header.Add(HeaderAppAuthToken, aat)
	} else if a.AppAuthToken != "" {
		req.Header.Add(HeaderAppAuthToken, a.AppAuthToken)
	}
	req.Header.Add("Accept", "application/json")
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Url: %s", url)
		a.logger.Debugf("Alipay_V3_Req_Body: %s", bm.JsonBody())
		a.logger.Debugf("Alipay_V3_Req_Headers: %#v", req.Header)
	}
	res, bs, err = req.Post(url).SendMultipartBodyMap(bm).EndBytesForAlipayV3(ctx)
	if err != nil {
		return nil, nil, err
	}
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Response: %d >> %s", res.StatusCode, string(bs))
		a.logger.Debugf("Alipay_V3_Rsp_Headers: %#v", res.Header)
	}
	return res, bs, nil
}

func (a *ClientV3) doDelete(ctx context.Context, bm gopay.BodyMap, uri, aat string) (res *http.Response, bs []byte, err error) {
	authorization, err := a.authorization(MethodDelete, uri, bm, aat, "")
	if err != nil {
		return nil, nil, err
	}

	var url = v3BaseUrlCh + uri
	if !a.IsProd {
		url = v3SandboxBaseUrl + uri
	}
	if a.proxyHost != "" {
		url = a.proxyHost + uri
	}
	req := a.hc.Req() // default json
	req.Header.Add(HeaderAuthorization, authorization)
	req.Header.Add(HeaderRequestID, a.requestIdFunc.RequestId())
	req.Header.Add(HeaderSdkVersion, "gopay/"+gopay.Version)
	if aat != gopay.NULL {
		req.Header.Add(HeaderAppAuthToken, aat)
	} else if a.AppAuthToken != "" {
		req.Header.Add(HeaderAppAuthToken, a.AppAuthToken)
	}
	req.Header.Add("Accept", "application/json")
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Url: %s", url)
		a.logger.Debugf("Alipay_V3_Req_Body: %s", bm.JsonBody())
		a.logger.Debugf("Alipay_V3_Req_Headers: %#v", req.Header)
	}
	res, bs, err = req.Delete(url).SendBodyMap(bm).EndBytesForAlipayV3(ctx)
	if err != nil {
		return nil, nil, err
	}

	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Response: %d >> %s", res.StatusCode, string(bs))
		a.logger.Debugf("Alipay_V3_Rsp_Headers: %#v", res.Header)
	}
	return res, bs, nil
}

func (a *ClientV3) encryptBizContent(originData string) (string, error) {
	encryptData, err := aes.CBCEncrypt([]byte(originData), []byte(a.aesKey), a.ivKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptData), nil
}

// encryptV3Body V3 接口内容加密
// 使用 AES-128-CBC 模式，IV 为 16 字节零向量
// 加密流程：JSON 明文 → AES/CBC/PKCS7 加密 → Base64 编码
func (a *ClientV3) encryptV3Body(plainText string) (string, error) {
	if plainText == "" {
		return "", nil
	}
	// go-pay/crypto/aes 的 CBCEncrypt 内部已处理 PKCS7 填充
	// 注意：aesKey 是 Base64 编码的，但 go-pay/crypto 的 CBCEncrypt 接受原始 key
	// 需要先把 aesKey Base64 解码为原始字节
	aesKeyBytes, err := base64.StdEncoding.DecodeString(a.aesKey)
	if err != nil {
		return "", fmt.Errorf("decode aes key error: %w", err)
	}
	encryptData, err := aes.CBCEncrypt([]byte(plainText), aesKeyBytes, a.ivKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptData), nil
}

// decryptV3Body V3 接口内容解密
// 解密流程：Base64 解码 → AES/CBC/PKCS7 解密 → JSON 明文
func (a *ClientV3) decryptV3Body(cipherText string) (string, error) {
	if cipherText == "" {
		return "", nil
	}
	cipherData, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", fmt.Errorf("decode cipher text error: %w", err)
	}
	aesKeyBytes, err := base64.StdEncoding.DecodeString(a.aesKey)
	if err != nil {
		return "", fmt.Errorf("decode aes key error: %w", err)
	}
	plainData, err := aes.CBCDecrypt(cipherData, aesKeyBytes, a.ivKey)
	if err != nil {
		return "", fmt.Errorf("aes cbc decrypt error: %w", err)
	}
	return string(plainData), nil
}
