package repository

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"

	"github.com/rezaAmiri123/service-article/internal/model"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *model.Article) error
	CreateComment(ctx context.Context, comment *model.Comment) error
}

type ORMArticleRepository struct {
	db *gorm.DB
}

func NewORMArticleRepository(db *gorm.DB)*ORMArticleRepository{
	return &ORMArticleRepository{db: db}
}

func (repo *ORMArticleRepository)Create(ctx context.Context, article *model.Article) error{
	span, ctx := opentracing.StartSpanFromContext(ctx, "ArticleRepository.Create")
	defer span.Finish()

	return repo.db.Create(article).Error
}

func (repo *ORMArticleRepository)CreateComment(ctx context.Context, comment *model.Comment) error{
	span, ctx := opentracing.StartSpanFromContext(ctx, "ArticleRepository.CreateComment")
	defer span.Finish()

	return repo.db.Create(comment).Error
}

