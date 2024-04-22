package article

import (
	"context"

	"gorm.io/gorm"
)

type ReaderDAO interface {
	Upsert(ctx context.Context, art Article) error
	UpsertV2(ctx context.Context, art PublishedArticle) error
}

func NewReaderDAO(db *gorm.DB) ReaderDAO {
	panic("implement me")
}
