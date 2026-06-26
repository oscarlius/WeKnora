package router

import (
	"testing"

	"github.com/Tencent/WeKnora/internal/types"
)

func TestAsynqQueueWeightsUsesEnvOverrides(t *testing.T) {
	t.Setenv("WEKNORA_ASYNQ_QUEUE_DEFAULT_WEIGHT", "2")
	t.Setenv("WEKNORA_ASYNQ_QUEUE_MULTIMODAL_WEIGHT", "4")

	weights := asynqQueueWeights()
	if got := weights[types.QueueDefault]; got != 2 {
		t.Fatalf("default queue weight = %d, want 2", got)
	}
	if got := weights[types.QueueMultimodal]; got != 4 {
		t.Fatalf("multimodal queue weight = %d, want 4", got)
	}
}

func TestAsynqQueueWeightsFallbackForInvalidValues(t *testing.T) {
	t.Setenv("WEKNORA_ASYNQ_QUEUE_DEFAULT_WEIGHT", "0")
	t.Setenv("WEKNORA_ASYNQ_QUEUE_MULTIMODAL_WEIGHT", "bad")

	weights := asynqQueueWeights()
	if got := weights[types.QueueDefault]; got != 3 {
		t.Fatalf("default queue weight = %d, want fallback 3", got)
	}
	if got := weights[types.QueueMultimodal]; got != 1 {
		t.Fatalf("multimodal queue weight = %d, want fallback 1", got)
	}
}
