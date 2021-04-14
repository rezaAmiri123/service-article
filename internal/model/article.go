package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"

	pb "github.com/rezaAmiri123/service-article/gen/pb"
)

// Article model
type Article struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Slug        string `gorm:"not null"`
	Description string `gorm:"not null"`
	Body        string `gorm:"not null"`
	Tags        []Tag  `gorm:"many2many:article_tags"`
	UserID      string `gorm:"not null"`
	Comments    []Comment
}

// Validate validates fields of article model
func (a Article) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Title, validation.Required),
		validation.Field(&a.Body, validation.Required),
		validation.Field(&a.Tags, validation.Required),
	)
}

// Overwrite overwrite each field if it's not zero-value
func (a *Article) Overwrite(title, description, body string) {
	if title != "" {
		a.Title = title
	}
	if description != "" {
		a.Description = description
	}
	if body != "" {
		a.Body = body
	}
}

// ProtoArticle generates proto article model from article
func (a *Article) ProtoArticle() *pb.Article {
	pa := pb.Article{
		Slug:        a.Slug,
		Title:       a.Title,
		Description: a.Description,
		Body:        a.Body,
	}

	// article tags
	tags := make([]string, 0, len(a.Tags))
	for _, t := range a.Tags {
		tags = append(tags, t.Name)
	}
	pa.TagList = tags
	return &pa
}
