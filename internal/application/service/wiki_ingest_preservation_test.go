package service

import (
	"testing"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/stretchr/testify/require"
)

func TestWikiBatchContextPreservesPageLimits(t *testing.T) {
	svc := &wikiIngestService{}
	ctx := svc.newWikiBatchContext("kb-1", &types.WikiConfig{
		MaxPageContentBytes: 4096,
		MaxRefs:             25,
	})

	require.Equal(t, 4096, ctx.MaxPageContentBytes)
	require.Equal(t, 25, ctx.MaxRefs)
}

func TestShouldSkipWikiPageResynthesis(t *testing.T) {
	page := &types.WikiPage{Content: "12345678"}
	ctx := &WikiBatchContext{MaxPageContentBytes: 8}
	additions := []SlugUpdate{{Type: types.WikiPageTypeEntity}}
	retracts := []SlugUpdate{{Type: "retract"}}

	require.True(t, shouldSkipWikiPageResynthesis(page, true, additions, nil, ctx))
	require.False(t, shouldSkipWikiPageResynthesis(page, false, additions, nil, ctx))
	require.False(t, shouldSkipWikiPageResynthesis(page, true, additions, retracts, ctx))
	require.False(t, shouldSkipWikiPageResynthesis(page, true, nil, nil, ctx))
	require.False(t, shouldSkipWikiPageResynthesis(page, true, additions, nil, &WikiBatchContext{}))
}

func TestCapRecentStringArray(t *testing.T) {
	values := types.StringArray{"old-1", "old-2", "new-1", "new-2"}
	require.Equal(t, types.StringArray{"new-1", "new-2"}, capRecentStringArray(values, 2))
	require.Equal(t, values, capRecentStringArray(values, 0))
}

func TestShouldSpawnKnowledgeSummary(t *testing.T) {
	require.True(t, shouldSpawnKnowledgeSummary(1, "summary-model"))
	require.False(t, shouldSpawnKnowledgeSummary(0, "summary-model"))
	require.False(t, shouldSpawnKnowledgeSummary(1, "  "))
}
