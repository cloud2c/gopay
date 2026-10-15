package xhttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

// 错误类别，供调用方用 errors.Is 判断，而不是解析错误文本。
var (
	// ErrNetwork 表示连接层失败（连接被拒、DNS、连接中断等）。
	ErrNetwork = errors.New("xhttp: network error")
	// ErrTimeout 表示请求超时（HTTP 客户端超时或上下文截止）。
	ErrTimeout = errors.New("xhttp: timeout")
	// ErrHTTPStatus 表示服务端返回了非预期的 HTTP 状态码（由调用方按需判断）。
	ErrHTTPStatus = errors.New("xhttp: http status")
)

// Error 是 xhttp 产生的传输层错误。错误信息只包含方法、主机、路径与类别，
// 不包含查询串、请求体与认证头。
type Error struct {
	Kind   error  // ErrNetwork、ErrTimeout 或 ErrHTTPStatus
	Method string // HTTP 方法
	Host   string // 主机（不含查询串）
	Path   string // 路径（不含查询串）
	Status int    // HTTP 状态码（仅 ErrHTTPStatus 时有效）
	Cause  error  // 底层错误，仅用于 Unwrap，不出现在 Error() 中
}

func (e *Error) Error() string {
	if e.Kind == ErrHTTPStatus {
		return fmt.Sprintf("%s %s%s: %s (%d)", e.Method, e.Host, e.Path, e.Kind, e.Status)
	}
	return fmt.Sprintf("%s %s%s: %s", e.Method, e.Host, e.Path, e.Kind)
}

// Unwrap 同时支持 errors.Is 的类别判断与底层错误的追溯。
func (e *Error) Unwrap() []error {
	if e.Cause == nil {
		return []error{e.Kind}
	}
	return []error{e.Kind, e.Cause}
}

// RetryPolicy 描述重试策略。零值表示不重试。
type RetryPolicy struct {
	MaxAttempts int           // 总尝试次数（含首次），≤1 表示不重试
	BaseDelay   time.Duration // 首次重试前等待，之后按 2 的幂增长
	MaxDelay    time.Duration // 单次等待上限
}

// SetRetry 设置重试策略。只对安全请求重试：GET，或携带 Idempotency-Key 的请求；
// 只对网络错误与 502、503、504 重试。POST 创建类请求不重试，避免重复扣款。
func (c *Client) SetRetry(p RetryPolicy) *Client {
	c.retry = p
	return c
}

// Event 是每次请求尝试结束时交给观测钩子的信息。
type Event struct {
	Method   string
	Host     string
	Path     string
	Status   int           // 0 表示没有收到响应
	Attempt  int           // 从 1 开始
	Duration time.Duration // 单次尝试耗时
	Err      error         // 传输层错误（*Error）或 nil
}

// Observer 在每次请求尝试结束后被调用。实现不得阻塞，也不得记录请求体或认证信息。
type Observer func(ctx context.Context, e Event)

// SetObserver 设置观测钩子，默认为空（不观测）。
func (c *Client) SetObserver(o Observer) *Client {
	c.observer = o
	return c
}

// classify 把底层错误归类为 *Error；无错误返回 nil。
func classify(method string, u *url.URL, err error) *Error {
	if err == nil {
		return nil
	}
	e := &Error{Method: method, Cause: err}
	if u != nil {
		e.Host = u.Host
		e.Path = u.Path
	}
	var netErr net.Error
	switch {
	case errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &netErr) && netErr.Timeout()):
		e.Kind = ErrTimeout
	default:
		e.Kind = ErrNetwork
	}
	return e
}

// retryable 判断某次尝试的结果是否值得重试（仅对安全请求调用）。
func retryable(status int, err error) bool {
	if err != nil {
		var e *Error
		return errors.As(err, &e) && e.Kind == ErrNetwork
	}
	return status == http.StatusBadGateway || status == http.StatusServiceUnavailable || status == http.StatusGatewayTimeout
}

// safeToRetry 判断请求是否可以安全重试。
func safeToRetry(method string, h http.Header) bool {
	return method == http.MethodGet || h.Get("Idempotency-Key") != ""
}

// backoff 返回第 attempt 次（从 1 开始）重试前的等待时长。
func (p RetryPolicy) backoff(attempt int) time.Duration {
	d := p.BaseDelay
	for i := 1; i < attempt; i++ {
		d *= 2
		if p.MaxDelay > 0 && d >= p.MaxDelay {
			return p.MaxDelay
		}
	}
	if p.MaxDelay > 0 && d > p.MaxDelay {
		return p.MaxDelay
	}
	return d
}

// sleepCtx 等待 d，若上下文先结束则返回其错误。
func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
