package modelusage

import (
	"context"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
)

type Recorder interface {
	Record(ctx context.Context, event types.ModelUsageEvent)
}

type noopRecorder struct{}

func (noopRecorder) Record(context.Context, types.ModelUsageEvent) {}

var (
	currentMu sync.RWMutex
	current   Recorder = noopRecorder{}
)

func SetRecorder(rec Recorder) {
	currentMu.Lock()
	defer currentMu.Unlock()
	if rec == nil {
		current = noopRecorder{}
		return
	}
	current = rec
}

func Record(ctx context.Context, event types.ModelUsageEvent) {
	currentMu.RLock()
	rec := current
	currentMu.RUnlock()
	rec.Record(ctx, event)
}

type AsyncRecorder struct {
	repo EventStore
	ch   chan types.ModelUsageEvent
	done chan struct{}
	wg   sync.WaitGroup
	once sync.Once
}

type EventStore interface {
	Create(ctx context.Context, event *types.ModelUsageEvent) error
}

func NewAsyncRecorder(repo EventStore) *AsyncRecorder {
	rec := &AsyncRecorder{
		repo: repo,
		ch:   make(chan types.ModelUsageEvent, 2048),
		done: make(chan struct{}),
	}
	rec.wg.Add(1)
	go rec.run()
	return rec
}

func (r *AsyncRecorder) Record(ctx context.Context, event types.ModelUsageEvent) {
	if r == nil || r.repo == nil {
		return
	}
	tenantID, ok := types.SessionTenantIDFromContext(ctx)
	if !ok || tenantID == 0 {
		return
	}
	event.TenantID = tenantID
	if userID, ok := types.UserIDFromContext(ctx); ok {
		event.UserID = userID
	}
	if requestID, ok := types.RequestIDFromContext(ctx); ok {
		event.RequestID = requestID
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}
	if event.TotalTokens == 0 {
		event.TotalTokens = event.PromptTokens + event.CompletionTokens
	}
	if event.UsageSource == "" {
		event.UsageSource = types.ModelUsageSourceMissing
	}

	select {
	case r.ch <- event:
	default:
		logger.Warnf(ctx, "model usage recorder queue full; dropping event model_id=%s kind=%s", event.ModelID, event.RequestKind)
	}
}

func (r *AsyncRecorder) Shutdown(ctx context.Context) error {
	if r == nil {
		return nil
	}
	r.once.Do(func() {
		close(r.ch)
		go func() {
			r.wg.Wait()
			close(r.done)
		}()
	})
	select {
	case <-r.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *AsyncRecorder) run() {
	defer r.wg.Done()
	for event := range r.ch {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		if err := r.repo.Create(ctx, &event); err != nil {
			logger.Warnf(ctx, "failed to persist model usage event: %v", err)
		}
		cancel()
	}
}

func ErrorType(err error) string {
	if err == nil {
		return ""
	}
	name := reflect.TypeOf(err).String()
	name = strings.TrimPrefix(name, "*")
	if name == "" {
		name = "error"
	}
	if len(name) > 128 {
		return name[:128]
	}
	return name
}
