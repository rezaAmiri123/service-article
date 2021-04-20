package model

import "github.com/jinzhu/gorm"

// Comment model
type FavoriteArticle struct {
	gorm.Model
	UserID    string   `gorm:"not null"`
	ArticleID uint   `gorm:"not null"`
	Article   Article
}
