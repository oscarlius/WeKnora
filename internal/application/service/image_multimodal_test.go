package service

import (
	"context"
	"errors"
	"testing"
)

type fakeImageVLM struct {
	id       string
	name     string
	response string
	err      error
	calls    int
}

func (f *fakeImageVLM) Predict(context.Context, [][]byte, string) (string, error) {
	f.calls++
	return f.response, f.err
}

func (f *fakeImageVLM) GetModelName() string { return f.name }
func (f *fakeImageVLM) GetModelID() string   { return f.id }

func TestPredictWithVLMFallback_PrimarySuccessDoesNotCallFallback(t *testing.T) {
	primary := &fakeImageVLM{id: "primary-id", name: "primary", response: "primary text"}
	fallback := &fakeImageVLM{id: "fallback-id", name: "fallback", response: "fallback text"}

	text, trace, err := predictWithVLMFallback(context.Background(), &resolvedImageVLM{
		primary:         primary,
		fallback:        fallback,
		primaryModelID:  "primary-id",
		fallbackModelID: "fallback-id",
	}, [][]byte{{1, 2, 3}}, "prompt")
	if err != nil {
		t.Fatalf("predict returned error: %v", err)
	}
	if text != "primary text" {
		t.Fatalf("expected primary response, got %q", text)
	}
	if primary.calls != 1 {
		t.Fatalf("expected primary called once, got %d", primary.calls)
	}
	if fallback.calls != 0 {
		t.Fatalf("expected fallback not called, got %d", fallback.calls)
	}
	if trace.UsedModelID != "primary-id" {
		t.Fatalf("expected used model primary-id, got %q", trace.UsedModelID)
	}
	if trace.FallbackReason != "" {
		t.Fatalf("expected no fallback reason, got %q", trace.FallbackReason)
	}
}

func TestPredictWithVLMFallback_PrimaryErrorCallsFallback(t *testing.T) {
	primaryErr := errors.New("upstream timeout")
	primary := &fakeImageVLM{id: "primary-id", name: "primary", err: primaryErr}
	fallback := &fakeImageVLM{id: "fallback-id", name: "fallback", response: "fallback text"}

	text, trace, err := predictWithVLMFallback(context.Background(), &resolvedImageVLM{
		primary:         primary,
		fallback:        fallback,
		primaryModelID:  "primary-id",
		fallbackModelID: "fallback-id",
	}, [][]byte{{1, 2, 3}}, "prompt")
	if err != nil {
		t.Fatalf("predict returned error: %v", err)
	}
	if text != "fallback text" {
		t.Fatalf("expected fallback response, got %q", text)
	}
	if primary.calls != 1 || fallback.calls != 1 {
		t.Fatalf("expected primary and fallback called once, got primary=%d fallback=%d", primary.calls, fallback.calls)
	}
	if trace.UsedModelID != "fallback-id" {
		t.Fatalf("expected used model fallback-id, got %q", trace.UsedModelID)
	}
	if trace.FallbackReason != primaryErr.Error() {
		t.Fatalf("expected fallback reason %q, got %q", primaryErr.Error(), trace.FallbackReason)
	}
}

func TestPredictWithVLMFallback_PrimaryEmptyTextDoesNotCallFallback(t *testing.T) {
	primary := &fakeImageVLM{id: "primary-id", name: "primary", response: ""}
	fallback := &fakeImageVLM{id: "fallback-id", name: "fallback", response: "fallback text"}

	text, trace, err := predictWithVLMFallback(context.Background(), &resolvedImageVLM{
		primary:         primary,
		fallback:        fallback,
		primaryModelID:  "primary-id",
		fallbackModelID: "fallback-id",
	}, [][]byte{{1, 2, 3}}, "prompt")
	if err != nil {
		t.Fatalf("predict returned error: %v", err)
	}
	if text != "" {
		t.Fatalf("expected empty primary response, got %q", text)
	}
	if fallback.calls != 0 {
		t.Fatalf("expected fallback not called for empty text, got %d", fallback.calls)
	}
	if trace.UsedModelID != "primary-id" {
		t.Fatalf("expected used model primary-id, got %q", trace.UsedModelID)
	}
}
