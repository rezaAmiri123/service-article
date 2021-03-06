package mysql

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/rezaAmiri123/service-article/cmd/config"
)

func NewGormDB(cfg *config.Config) *gorm.DB {
	DBString := "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
	URL := fmt.Sprintf(DBString, cfg.Database.DBUser, cfg.Database.DBPass, cfg.Database.DBHost, cfg.Database.DBPort, cfg.Database.DBName)
	db, err := gorm.Open(cfg.Database.DBType, URL)
	if err != nil {
		log.Fatal("cannot connect to the database", err)
	}
	return db
}
