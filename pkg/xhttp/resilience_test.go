package xhttp

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// 测试用的快速重试策略，避免测试等待真实退避时间。
var fastRetry = RetryPolicy{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}

func TestRetryGETOnServiceUnavailable(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if hits.Add(1) < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := NewClient().SetRetry(fastRetry)
	res, bs, err := c.Req().Get(srv.URL + "/x").EndBytes(context.Background())
	if err != nil || res.StatusCode != 200 || string(bs) != `{"ok":true}` {
		t.Fatalf("want success after retries, got status=%v body=%q err=%v", res, bs, err)
	}
	if hits.Load() != 3 {
		t.Fatalf("hits = %d, want 3", hits.Load())
	}
}

func TestNoRetryForUnsafePOSTWithoutKey(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	c := NewClient().SetRetry(fastRetry)
	_, _, _ = c.Req(TypeJSON).Post(srv.URL + "/pay").SendStruct(map[string]string{"a": "b"}).EndBytes(context.Background())
	if hits.Load() != 1 {
		t.Fatalf("POST without Idempotency-Key was retried: hits = %d", hits.Load())
	}
}

func TestRetryPOSTWithIdempotencyKey(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if hits.Add(1) == 1 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		_, _ = w.Write([]byte("done"))
	}))
	defer srv.Close()

	c := NewClient().SetRetry(fastRetry)
	req := c.Req(TypeJSON).Post(srv.URL + "/pay")
	req.Header.Set("Idempotency-Key", "k-1")
	_, bs, err := req.SendStruct(map[string]string{"a": "b"}).EndBytes(context.Background())
	if err != nil || string(bs) != "done" || hits.Load() != 2 {
		t.Fatalf("keyed POST: body=%q err=%v hits=%d", bs, err, hits.Load())
	}
}

func TestNoRetryOn500(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient().SetRetry(fastRetry)
	_, _, _ = c.Req().Get(srv.URL).EndBytes(context.Background())
	if hits.Load() != 1 {
		t.Fatalf("500 was retried: hits = %d", hits.Load())
	}
}

func TestNetworkErrorIsClassifiedAndRedacted(t *testing.T) {
	c := NewClient()
	// 使用已关闭的端口制造连接失败，查询串中放入敏感值。
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	addr := srv.URL
	srv.Close()

	_, _, err := c.Req().Get(addr + "/path?secret=TOPSECRET").EndBytes(context.Background())
	if err == nil {
		t.Fatal("want network error")
	}
	if !errors.Is(err, ErrNetwork) {
		t.Fatalf("want ErrNetwork, got %v", err)
	}
	if strings.Contains(err.Error(), "TOPSECRET") || strings.Contains(err.Error(), "secret=") {
		t.Fatalf("error leaks query string: %q", err.Error())
	}
}

func TestTimeoutIsClassified(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()

	c := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, _, err := c.Req().Get(srv.URL).EndBytes(ctx)
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("want ErrTimeout, got %v", err)
	}
}

func TestObserverSeesEachAttemptWithoutQueryOrBody(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if hits.Add(1) < 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	var events []Event
	c := NewClient().SetRetry(fastRetry).SetObserver(func(_ context.Context, e Event) {
		events = append(events, e)
	})
	_, _, _ = c.Req().Get(srv.URL + "/v1/orders?token=SECRET").EndBytes(context.Background())
	if len(events) != 2 {
		t.Fatalf("events = %d, want 2", len(events))
	}
	if events[0].Attempt != 1 || events[0].Status != 503 || events[1].Attempt != 2 || events[1].Status != 200 {
		t.Fatalf("unexpected events: %+v", events)
	}
	for _, e := range events {
		if strings.Contains(e.Path, "SECRET") || e.Host == "" || e.Path != "/v1/orders" {
			t.Fatalf("event leaks or misses fields: %+v", e)
		}
	}
}

func TestBackoffCapsAtMaxDelay(t *testing.T) {
	p := RetryPolicy{BaseDelay: 100 * time.Millisecond, MaxDelay: 300 * time.Millisecond}
	want := []time.Duration{100 * time.Millisecond, 200 * time.Millisecond, 300 * time.Millisecond, 300 * time.Millisecond}
	for i, w := range want {
		if got := p.backoff(i + 1); got != w {
			t.Fatalf("backoff(%d) = %v, want %v", i+1, got, w)
		}
	}
}
