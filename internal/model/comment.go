package model

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"

	pb "github.com/rezaAmiri123/service-article/gen/pb"
)

// Comment model
type Comment struct {
	gorm.Model
	Body      string `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	ArticleID uint   `gorm:"not null"`
	Article   Article
}

// Validate validates fields of comment model
func (c Comment) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Body, validation.Required),
	)
}

// ProtoComment generates proto comment model from article
func (c *Comment) ProtoComment() *pb.Comment {
	return &pb.Comment{
		Id:   fmt.Sprintf("%d", c.ID),
		Body: c.Body,
	}
}
