package repository

import (
	"context"
	"github.com/rezaAmiri123/service-article/pkg/utils"

	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"

	"github.com/rezaAmiri123/service-article/internal/model"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *model.Article) error
	Update(ctx context.Context, article *model.Article) error
	Delete(ctx context.Context, article *model.Article) error
	GetBySlug(ctx context.Context, slug string) (*model.Article, error)
	GetByID(ctx context.Context, id string) (*model.Article, error)
	CreateComment(ctx context.Context, comment *model.Comment) error
	GetCommentByID(ctx context.Context, id string) (*model.Comment, error)
	GetComments(ctx context.Context, article *model.Article) ([]model.Comment, error)
	DeleteComment(ctx context.Context, comment *model.Comment) error
}

type ORMArticleRepository struct {
	db *gorm.DB
}

func NewORMArticleRepository(db *gorm.DB) *ORMArticleRepository {
	return &ORMArticleRepository{db: db}
}

func (repo *ORMArticleRepository) Create(ctx context.Context, article *model.Article) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.Create")
	defer span.Finish()

	return repo.db.Create(article).Error
}

func (repo *ORMArticleRepository) CreateComment(ctx context.Context, comment *model.Comment) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.CreateComment")
	defer span.Finish()

	return repo.db.Create(comment).Error
}

func (repo *ORMArticleRepository) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.GetBySlug")
	defer span.Finish()

	var a model.Article
	if err := repo.db.Where(model.Article{Slug: slug}).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (repo *ORMArticleRepository) GetByID(ctx context.Context, id string) (*model.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.GetByID")
	defer span.Finish()

	var a model.Article
	if err := repo.db.Preload("Tags").Find(&a, utils.StringToUint(id)).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (repo *ORMArticleRepository) Update(ctx context.Context, article *model.Article) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.Update")
	defer span.Finish()

	return repo.db.Model(article).Update(article).Error
}

func (repo *ORMArticleRepository) Delete(ctx context.Context, article *model.Article) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.Delete")
	defer span.Finish()

	return repo.db.Delete(article).Error
}

func (repo *ORMArticleRepository) GetCommentByID(ctx context.Context, id string) (*model.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.GetCommentByID")
	defer span.Finish()

	var m model.Comment
	err := repo.db.Find(&m, utils.StringToUint(id)).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (repo *ORMArticleRepository)DeleteComment(ctx context.Context, comment *model.Comment) error{
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.DeleteComment")
	defer span.Finish()

	return repo.db.Delete(comment).Error
}

func (repo *ORMArticleRepository)GetComments(ctx context.Context, article *model.Article) ([]model.Comment, error){
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.GetComments")
	defer span.Finish()

	var cs []model.Comment
	err := repo.db.Where(model.Comment{
		ArticleID: article.ID,
	}).Find(&cs).Error
	if err !=nil{
		return nil, err
	}
	return cs, nil
}
