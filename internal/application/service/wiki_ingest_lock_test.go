package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
)

type wikiIngestLockPendingRepo struct {
	count int64
}

func (r *wikiIngestLockPendingRepo) Enqueue(context.Context, *types.TaskPendingOp) error {
	return nil
}

func (r *wikiIngestLockPendingRepo) PeekBatch(context.Context, string, string, string, int) ([]*types.TaskPendingOp, error) {
	return nil, nil
}

func (r *wikiIngestLockPendingRepo) DeleteByIDs(context.Context, []int64) error {
	return nil
}

func (r *wikiIngestLockPendingRepo) IncrFailCount(context.Context, int64) (int, error) {
	return 0, nil
}

func (r *wikiIngestLockPendingRepo) PendingCount(context.Context, string, string, string) (int64, error) {
	return r.count, nil
}

func (r *wikiIngestLockPendingRepo) DeleteByDedupKey(context.Context, string, string, string, string, string) error {
	return nil
}

type wikiIngestLockTaskEnqueuer struct {
	tasks []*asynq.Task
	err   error
}

func (e *wikiIngestLockTaskEnqueuer) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	e.tasks = append(e.tasks, task)
	if e.err != nil {
		return nil, e.err
	}
	return &asynq.TaskInfo{ID: "follow-up", Queue: "low"}, nil
}

func TestProcessWikiIngest_LiteLockConflictSchedulesFollowUpWithoutRetry(t *testing.T) {
	enqueuer := &wikiIngestLockTaskEnqueuer{}
	svc := &wikiIngestService{
		task:        enqueuer,
		pendingRepo: &wikiIngestLockPendingRepo{count: 3},
	}
	svc.liteLocks.Store("kb-1", struct{}{})

	raw, err := json.Marshal(WikiIngestPayload{
		TenantID:        10000,
		KnowledgeBaseID: "kb-1",
		Language:        "zh-CN",
	})
	require.NoError(t, err)

	err = svc.ProcessWikiIngest(context.Background(), asynq.NewTask(types.TypeWikiIngest, raw))
	require.NoError(t, err)
	require.Len(t, enqueuer.tasks, 1)
	require.Equal(t, types.TypeWikiIngest, enqueuer.tasks[0].Type())
}

func TestProcessWikiIngest_LiteLockConflictWithEmptyQueueReturnsNil(t *testing.T) {
	enqueuer := &wikiIngestLockTaskEnqueuer{}
	svc := &wikiIngestService{
		task:        enqueuer,
		pendingRepo: &wikiIngestLockPendingRepo{count: 0},
	}
	svc.liteLocks.Store("kb-1", struct{}{})

	raw, err := json.Marshal(WikiIngestPayload{
		TenantID:        10000,
		KnowledgeBaseID: "kb-1",
	})
	require.NoError(t, err)

	err = svc.ProcessWikiIngest(context.Background(), asynq.NewTask(types.TypeWikiIngest, raw))
	require.NoError(t, err)
	require.Empty(t, enqueuer.tasks)
}

func TestScheduleFollowUp_TaskIDConflictCountsAsScheduled(t *testing.T) {
	enqueuer := &wikiIngestLockTaskEnqueuer{err: asynq.ErrTaskIDConflict}
	svc := &wikiIngestService{
		task:        enqueuer,
		pendingRepo: &wikiIngestLockPendingRepo{count: 2},
	}

	ok := svc.scheduleFollowUp(context.Background(), WikiIngestPayload{
		TenantID:        10000,
		KnowledgeBaseID: "kb-1",
		Language:        "zh-CN",
	})

	require.True(t, ok)
	require.Len(t, enqueuer.tasks, 1)
	require.True(t, errors.Is(enqueuer.err, asynq.ErrTaskIDConflict))
}
