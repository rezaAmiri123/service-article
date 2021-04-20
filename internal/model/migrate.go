package model

import "github.com/jinzhu/gorm"

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&FavoriteArticle{},
		&Tag{},
		&Comment{},
		&Article{},
	).Error
}
