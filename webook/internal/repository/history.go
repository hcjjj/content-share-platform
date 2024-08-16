package repository

import (
	"basic-go/webook/internal/domain"
	"context"
)

type HistoryRecordRepository interface {
	AddRecord(ctx context.Context, record domain.HistoryRecord) error
}
