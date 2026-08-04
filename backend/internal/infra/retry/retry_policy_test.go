package retry

import (
	"context"
	"errors"
	"testing"
)

func TestRetryPolicy_Execute_Success(t *testing.T) {
	policy := NewRetryPolicy()

	calls := 0
	fn := func(ctx context.Context) error {
		calls++
		return nil
	}

	err := policy.Execute(context.Background(), "test", fn)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetryPolicy_Execute_RetrySuccess(t *testing.T) {
	policy := NewRetryPolicy()

	calls := 0
	fn := func(ctx context.Context) error {
		calls++
		if calls < 2 {
			return errors.New("temporary error")
		}
		return nil
	}

	err := policy.Execute(context.Background(), "test", fn)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestRetryPolicy_Execute_MaxAttempts(t *testing.T) {
	policy := NewRetryPolicy()

	calls := 0
	fn := func(ctx context.Context) error {
		calls++
		return errors.New("persistent error")
	}

	err := policy.Execute(context.Background(), "test", fn)
	if err == nil {
		t.Error("expected error")
	}

	if calls != policy.MaxAttempts {
		t.Errorf("expected %d calls, got %d", policy.MaxAttempts, calls)
	}
}

func TestRetryPolicy_Execute_NonRetryable(t *testing.T) {
	policy := NewRetryPolicy()
	policy.Retryable = func(err error) bool {
		return false
	}

	calls := 0
	fn := func(ctx context.Context) error {
		calls++
		return errors.New("non-retryable error")
	}

	err := policy.Execute(context.Background(), "test", fn)
	if err == nil {
		t.Error("expected error")
	}

	if calls != 1 {
		t.Errorf("expected 1 call for non-retryable error, got %d", calls)
	}
}

func TestRetryPolicy_Execute_ContextCancelled(t *testing.T) {
	policy := NewRetryPolicy()

	calls := 0
	fn := func(ctx context.Context) error {
		calls++
		return errors.New("error")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := policy.Execute(ctx, "test", fn)
	if err != context.Canceled {
		t.Errorf("expected context cancelled error, got %v", err)
	}
}

func TestRedisRetryPolicy(t *testing.T) {
	calls := 0
	fn := func(ctx context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("redis error")
		}
		return nil
	}

	err := RedisRetryPolicy.Execute(context.Background(), "redis", fn)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDatabaseRetryPolicy_NoRetry(t *testing.T) {
	calls := 0
	fn := func(ctx context.Context) error {
		calls++
		return errors.New("db error")
	}

	err := DatabaseRetryPolicy.Execute(context.Background(), "database", fn)
	if err == nil {
		t.Error("expected error")
	}

	if calls != 1 {
		t.Errorf("expected 1 call (no retry), got %d", calls)
	}
}
