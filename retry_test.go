package redis

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

func TestRetryManager_ctxCancelDuringOp(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	ctx, cancel := context.WithCancel(context.Background())
	err := rm.execute(ctx, func(context.Context) error { cancel(); return errors.New("x") })
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestRetryManager_nilCtxAndZeroSleep(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	//lint:ignore SA1012 intentional nil context in test
	if err := rm.execute(nil, func(context.Context) error { return nil }); err != nil {
		t.Fatal(err)
	}
	if err := rm.sleep(context.Background(), 0); err != nil {
		t.Fatal(err)
	}
}

func TestRetryManager_ctxCancelDuringBackoff(t *testing.T) {
	rm := newRetryManager(RetryConfig{MaxAttempts: 5, InitialDelay: time.Hour, MaxDelay: time.Hour, Multiplier: 1})
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(20 * time.Millisecond); cancel() }()
	if err := rm.execute(ctx, func(context.Context) error { return errors.New("x") }); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
