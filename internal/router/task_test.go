package router

import (
	"testing"

	"github.com/Tencent/WeKnora/internal/application/service"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/hibiken/asynq"
)

func TestAsynqRetryDelayUsesShortBackoffForVLMRateLimit(t *testing.T) {
	task := asynq.NewTask(types.TypeImageMultimodal, nil)
	if got := asynqRetryDelayFunc(1, service.ErrVLMRateLimited, task); got != vlmRateLimitRetryDelay {
		t.Fatalf("retry delay = %s, want %s", got, vlmRateLimitRetryDelay)
	}
}
