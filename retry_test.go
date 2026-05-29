package goredis

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryManager_execute_success(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	n := 0
	err := rm.execute(context.Background(), func(ctx context.Context) error {
		n++
		if n < 2 {
			return errors.New("retry")
		}
		return nil
	})
	if err != nil || n != 2 {
		t.Fatalf("err=%v n=%d", err, n)
	}
}

func TestRetryManager_execute_exhausted(t *testing.T) {
	rm := newRetryManager(RetryConfig{MaxAttempts: 2, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, Multiplier: 2})
	err := rm.execute(context.Background(), func(ctx context.Context) error {
		return errors.New("fail")
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRetryManager_contextCancel(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := rm.execute(ctx, func(ctx context.Context) error { return errors.New("x") })
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRetryManager_sleepError(t *testing.T) {
	rm := newRetryManager(RetryConfig{MaxAttempts: 2, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, Multiplier: 1})
	rm.sleep = func(context.Context, time.Duration) error { return errors.New("sleep") }
	err := rm.execute(context.Background(), func(context.Context) error { return errors.New("op") })
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_normalizeRetry(t *testing.T) {
	r := RetryConfig{}
	normalizeRetry(&r)
	if r.MaxAttempts != DefaultRetryConfig().MaxAttempts {
		t.Fatal(r)
	}
}

func TestRetryManager_jitter(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	if d := rm.jitter(0); d != 0 {
		t.Fatal(d)
	}
}
