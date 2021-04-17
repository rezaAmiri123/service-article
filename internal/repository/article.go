package repository

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"

	"github.com/rezaAmiri123/service-article/internal/model"
	"github.com/rezaAmiri123/service-article/pkg/utils"
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
	AddFavorite(ctx context.Context, article *model.Article, userID string) error
	DeleteFavorite(ctx context.Context, article *model.Article, userID string) error
	IsFavorited(ctx context.Context, article *model.Article, userID string) (bool, error)
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

func (repo *ORMArticleRepository) DeleteComment(ctx context.Context, comment *model.Comment) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.DeleteComment")
	defer span.Finish()

	return repo.db.Delete(comment).Error
}

func (repo *ORMArticleRepository) GetComments(ctx context.Context, article *model.Article) ([]model.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.GetComments")
	defer span.Finish()

	var cs []model.Comment
	err := repo.db.Where(model.Comment{
		ArticleID: article.ID,
	}).Find(&cs).Error
	if err != nil {
		return nil, err
	}
	return cs, nil
}

func (repo *ORMArticleRepository) AddFavorite(ctx context.Context, article *model.Article, userID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.AddFavorite")
	defer span.Finish()

	tx := repo.db.Begin()
	err := tx.Model(article).Association("FavoritedUsers").Append(userID).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(article).Update("favorites_count", gorm.Expr("favorites_count + ?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	article.FavoritesCount++
	return nil
}

func (repo *ORMArticleRepository) DeleteFavorite(ctx context.Context, article *model.Article, userID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.DeleteFavorite")
	defer span.Finish()

	tx := repo.db.Begin()
	err := tx.Model(article).Association("FavoritedUsers").Delete(userID).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(article).Update("favorites_count", gorm.Expr("favorites_count - ?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	article.FavoritesCount--

	return nil
}

func (repo *ORMArticleRepository) IsFavorited(ctx context.Context, article *model.Article, userID string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ORMArticleRepository.IsFavorited")
	defer span.Finish()

	if article == nil || userID == "" {
		return false, nil
	}

	var count int
	filter := model.Article{}
	filter.ID = article.ID
	filter.UserID = userID
	err := repo.db.Table("favorite_articles").Where(filter).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
