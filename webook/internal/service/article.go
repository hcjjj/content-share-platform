package service

import (
	"context"

	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
}

type ArticleServiceV1 struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &ArticleServiceV1{
		repo: repo,
	}
}

func (a *ArticleServiceV1) Save(ctx context.Context, art domain.Article) (int64, error) {
	return a.repo.Create(ctx, art)
}
