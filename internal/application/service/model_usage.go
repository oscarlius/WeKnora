package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

type modelUsageService struct {
	repo interfaces.ModelUsageRepository
}

func NewModelUsageService(repo interfaces.ModelUsageRepository) interfaces.ModelUsageService {
	return &modelUsageService{repo: repo}
}

func (s *modelUsageService) GetUsageReport(
	ctx context.Context,
	query types.ModelUsageQuery,
) (*types.ModelUsageReport, error) {
	tenantID, ok := types.TenantIDFromContext(ctx)
	if !ok || tenantID == 0 {
		return nil, fmt.Errorf("tenant ID not found in context")
	}

	now := time.Now().UTC()
	query.End = now
	query.Start, query.BucketSize = usageWindow(query.Range, now)
	query.RecentLimit = 30

	if query.ModelType == types.ModelType("all") {
		query.ModelType = ""
	}
	return s.repo.Report(ctx, tenantID, query)
}

func usageWindow(value string, now time.Time) (time.Time, time.Duration) {
	switch value {
	case "15m":
		return now.Add(-15 * time.Minute), time.Minute
	case "1h":
		return now.Add(-time.Hour), 5 * time.Minute
	case "7d":
		return now.Add(-7 * 24 * time.Hour), 6 * time.Hour
	default:
		return now.Add(-24 * time.Hour), 30 * time.Minute
	}
}
