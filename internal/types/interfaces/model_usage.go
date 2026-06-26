package interfaces

import (
	"context"

	"github.com/Tencent/WeKnora/internal/types"
)

type ModelUsageRepository interface {
	Create(ctx context.Context, event *types.ModelUsageEvent) error
	Report(ctx context.Context, tenantID uint64, query types.ModelUsageQuery) (*types.ModelUsageReport, error)
}

type ModelUsageService interface {
	GetUsageReport(ctx context.Context, query types.ModelUsageQuery) (*types.ModelUsageReport, error)
}
