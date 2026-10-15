package xhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-pay/xlog"

	"github.com/cloud2c/gopay"
)

type RequestIdHandler interface {
	RequestId() string
}

type Request struct {
	client           *Client
	Header           http.Header
	formString       string
	jsonByte         []byte
	url              string
	method           string
	requestType      string
	responseType     string
	multipartBodyMap map[string]any
	err              error
}

func (r *Request) Get(url string) *Request {
	r.method = GET
	r.url = url
	return r
}

func (r *Request) Post(url string) *Request {
	r.method = POST
	r.url = url
	return r
}

func (r *Request) Put(url string) *Request {
	r.method = PUT
	r.url = url
	return r
}

func (r *Request) Delete(url string) *Request {
	r.method = DELETE
	r.url = url
	return r
}

func (r *Request) Patch(url string) *Request {
	r.method = PATCH
	r.url = url
	return r
}

// =====================================================================================================================

func (r *Request) SendStruct(v any) (c *Request) {
	if v == nil {
		return r
	}
	bs, err := json.Marshal(v)
	if err != nil {
		r.err = fmt.Errorf("json.Marshal(%+v): %w", v, err)
		return r
	}
	switch r.requestType {
	case TypeJSON:
		r.jsonByte = bs
	case TypeXML, TypeFormData:
		body := make(map[string]any)
		if err = json.Unmarshal(bs, &body); err != nil {
			r.err = fmt.Errorf("json.Unmarshal(%s, %+v): %w", string(bs), body, err)
			return r
		}
		r.formString = FormatURLParam(body)
	}
	return r
}

func (r *Request) SendBodyMap(bm map[string]any) (client *Request) {
	if bm == nil {
		return r
	}
	switch r.requestType {
	case TypeJSON:
		bs, err := json.Marshal(bm)
		if err != nil {
			r.err = fmt.Errorf("json.Marshal(%+v): %w", bm, err)
			return r
		}
		r.jsonByte = bs
	case TypeXML, TypeFormData:
		r.formString = FormatURLParam(bm)
	}
	return r
}

func (r *Request) SendMultipartBodyMap(bm map[string]any) (client *Request) {
	if bm == nil {
		return r
	}
	switch r.requestType {
	case TypeJSON:
		bs, err := json.Marshal(bm)
		if err != nil {
			r.err = fmt.Errorf("json.Marshal(%+v): %w", bm, err)
			return r
		}
		r.jsonByte = bs
	case TypeXML, TypeFormData:
		r.formString = FormatURLParam(bm)
	case TypeMultipartFormData:
		r.multipartBodyMap = bm
	}
	return r
}

// encodeStr: url.Values.Encode() or jsonBody
func (r *Request) SendString(encodeStr string) (client *Request) {
	switch r.requestType {
	case TypeJSON:
		r.jsonByte = []byte(encodeStr)
	case TypeXML, TypeFormData:
		r.formString = encodeStr
	}
	return r
}

// =====================================================================================================================

func (r *Request) EndBytesForAlipayV3(ctx context.Context) (res *http.Response, bs []byte, err error) {
	if r.err != nil {
		return nil, nil, r.err
	}
	var (
		body io.Reader
		bw   *multipart.Writer
	)
	// multipart-form-data
	if r.requestType == TypeMultipartFormData {
		body = &bytes.Buffer{}
		bw = multipart.NewWriter(body.(io.Writer))
	}

	switch r.method {
	case GET:
		// do nothing
	case POST, PUT, DELETE, PATCH:
		switch r.requestType {
		case TypeJSON:
			if r.jsonByte != nil {
				body = strings.NewReader(string(r.jsonByte))
			}
		case TypeFormData:
			if r.formString != "" {
				body = strings.NewReader(r.formString)
			}
		case TypeMultipartFormData:
			for k, v := range r.multipartBodyMap {
				// file_content 参数
				if file, ok := v.(*gopay.File); ok {
					fw, e := bw.CreateFormFile(k, file.Name)
					if e != nil {
						return nil, nil, e
					}
					_, _ = fw.Write(file.Content)

					//h := make(textproto.MIMEHeader)
					//h.Set("Content-Disposition",
					//	fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
					//		escapeQuotes(k), escapeQuotes(file.Name)))
					//h.Set("Content-Type", "application/octet-stream")
					//part, e := bw.CreatePart(h)
					//if e != nil {
					//	return nil, nil, e
					//}
					//_, _ = part.Write(file.Content)
					continue
				}

				// text 参数 data
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, escapeQuotes(k)))
				h.Set("Content-Type", "application/json")
				part, e := bw.CreatePart(h)
				if e != nil {
					return nil, nil, e
				}
				switch vs := v.(type) {
				case string:
					_, _ = part.Write([]byte(vs))
				default:
					xlog.Warnf("multipart-form-data key: %s, value: %s", k, ConvertToString(v))
					_, _ = part.Write([]byte(ConvertToString(v)))
				}
			}
			_ = bw.Close()
			r.Header.Set("Content-Type", bw.FormDataContentType())
		case TypeXML:
			if r.formString != "" {
				body = strings.NewReader(r.formString)
			}
		default:
			return nil, nil, errors.New("Request type Error ")
		}
	default:
		return nil, nil, errors.New("Only support GET and POST and PUT and DELETE ")
	}

	// request
	req, err := http.NewRequestWithContext(ctx, r.method, r.url, body)
	if err != nil {
		return nil, nil, err
	}
	req.Header = r.Header
	res, err = r.client.HttpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	bs, err = io.ReadAll(io.LimitReader(res.Body, int64(r.client.bodySizeMB<<20))) // default 10MB change the size you want
	if err != nil {
		return nil, nil, err
	}
	return res, bs, nil
}

func (r *Request) EndBytes(ctx context.Context) (res *http.Response, bs []byte, err error) {
	if r.err != nil {
		return nil, nil, r.err
	}
	payload, err := r.buildPayload()
	if err != nil {
		return nil, nil, err
	}
	u, _ := url.Parse(r.url)
	attempts := 1
	if safeToRetry(r.method, r.Header) && r.client.retry.MaxAttempts > 1 {
		attempts = r.client.retry.MaxAttempts
	}
	for attempt := 1; ; attempt++ {
		res, bs, err = r.doOnce(ctx, u, payload, attempt)
		status := 0
		if res != nil {
			status = res.StatusCode
		}
		if attempt >= attempts || !retryable(status, err) {
			return res, bs, err
		}
		if e := sleepCtx(ctx, r.client.retry.backoff(attempt)); e != nil {
			return nil, nil, e
		}
	}
}

// buildPayload 按请求类型生成本次请求的请求体字节；GET 与无请求体的请求返回 nil。
// 每次尝试使用新的 Reader，因此可以安全重试。
func (r *Request) buildPayload() ([]byte, error) {
	switch r.method {
	case GET:
		return nil, nil
	case POST, PUT, DELETE, PATCH:
	default:
		return nil, errors.New("Only support GET and POST and PUT and DELETE ")
	}
	switch r.requestType {
	case TypeJSON:
		if r.jsonByte != nil {
			return r.jsonByte, nil
		}
		return nil, nil
	case TypeFormData, TypeXML:
		if r.formString != "" {
			return []byte(r.formString), nil
		}
		return nil, nil
	case TypeMultipartFormData:
		buf := &bytes.Buffer{}
		bw := multipart.NewWriter(buf)
		for k, v := range r.multipartBodyMap {
			// file 参数
			if file, ok := v.(*gopay.File); ok {
				fw, e := bw.CreateFormFile(k, file.Name)
				if e != nil {
					return nil, e
				}
				_, _ = fw.Write(file.Content)
				continue
			}
			// text 参数
			switch vs := v.(type) {
			case string:
				_ = bw.WriteField(k, vs)
			default:
				_ = bw.WriteField(k, ConvertToString(v))
			}
		}
		_ = bw.Close()
		r.Header.Set("Content-Type", bw.FormDataContentType())
		return buf.Bytes(), nil
	default:
		return nil, errors.New("Request type Error ")
	}
}

// doOnce 执行一次请求尝试，并通知观测钩子。
func (r *Request) doOnce(ctx context.Context, u *url.URL, payload []byte, attempt int) (*http.Response, []byte, error) {
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, r.method, r.url, body)
	if err != nil {
		return nil, nil, err
	}
	req.Header = r.Header
	res, err := r.client.HttpClient.Do(req)
	var bs []byte
	if err == nil {
		defer res.Body.Close()
		bs, err = io.ReadAll(io.LimitReader(res.Body, int64(r.client.bodySizeMB<<20))) // default 10MB change the size you want
	}
	if r.client.observer != nil {
		status := 0
		if res != nil {
			status = res.StatusCode
		}
		ev := Event{Method: r.method, Attempt: attempt, Duration: time.Since(start), Status: status}
		if u != nil {
			ev.Host = u.Host
			ev.Path = u.Path
		}
		if err != nil {
			ev.Err = classify(r.method, u, err)
		}
		r.client.observer(ctx, ev)
	}
	if err != nil {
		return nil, nil, classify(r.method, u, err)
	}
	return res, bs, nil
}

func (r *Request) EndStruct(ctx context.Context, v any) (res *http.Response, err error) {
	res, bs, err := r.EndBytes(ctx)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return res, fmt.Errorf("StatusCode(%d) != 200", res.StatusCode)
	}

	switch r.responseType {
	case ResTypeJSON:
		err = json.Unmarshal(bs, &v)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal(%s, %+v): %w", string(bs), v, err)
		}
		return res, nil
	case ResTypeXML:
		err = xml.Unmarshal(bs, &v)
		if err != nil {
			return nil, fmt.Errorf("xml.Unmarshal(%s, %+v): %w", string(bs), v, err)
		}
		return res, nil
	default:
		return nil, errors.New("responseType Type Wrong")
	}
}

func FormatURLParam(body map[string]any) (urlParam string) {
	var (
		buf  strings.Builder
		keys []string
	)
	for k := range body {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v, ok := body[k].(string)
		if !ok {
			v = ConvertToString(body[k])
		}
		if v != "" {
			buf.WriteString(url.QueryEscape(k))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
			buf.WriteByte('&')
		}
	}
	if buf.Len() <= 0 {
		return ""
	}
	return buf.String()[:buf.Len()-1]
}
