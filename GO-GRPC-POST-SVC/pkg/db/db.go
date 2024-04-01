package db

import (
	"fmt"

	"github.com/akhi9550/post-svc/pkg/config"
	"github.com/akhi9550/post-svc/pkg/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	db, dbErr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	db.AutoMigrate(&domain.Post{})
	db.AutoMigrate(&domain.Url{})
	db.AutoMigrate(&domain.PostType{})
	db.AutoMigrate(&domain.Likes{})
	db.AutoMigrate(&domain.Comment{})
	db.AutoMigrate(&domain.CommentRepies{})
	db.AutoMigrate(&domain.Tags{})
	db.AutoMigrate(&domain.Story{})
	db.AutoMigrate(&domain.StoryLike{})
	db.AutoMigrate(&domain.ViewStory{})
	db.AutoMigrate(&domain.SavedPost{})
	return db, dbErr
}
