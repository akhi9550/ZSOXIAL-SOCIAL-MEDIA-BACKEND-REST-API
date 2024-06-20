package db

import (
	"fmt"

	"github.com/akhi9550/post-svc/pkg/config"
	"github.com/akhi9550/post-svc/pkg/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=postgres port=%s password=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPort, cfg.DBPassword)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	rows, err := db.Raw("SELECT 1 FROM pg_database WHERE datname = ?", cfg.DBName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
		if err := db.Exec(createDBQuery).Error; err != nil {
			return nil, err
		}
		fmt.Printf("Database %s created successfully.\n", cfg.DBName)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Close(); err != nil {
		return nil, err
	}

	psqlInfo = fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	db, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&domain.Post{})
	db.AutoMigrate(&domain.PostType{})
	db.AutoMigrate(&domain.Likes{})
	db.AutoMigrate(&domain.Comment{})
	db.AutoMigrate(&domain.CommentReplies{})
	db.AutoMigrate(&domain.ArchivePost{})
	db.AutoMigrate(&domain.Tags{})
	db.AutoMigrate(&domain.Story{})
	db.AutoMigrate(&domain.StoryLike{})
	db.AutoMigrate(&domain.ArchiveStory{})
	db.AutoMigrate(&domain.ViewStory{})
	db.AutoMigrate(&domain.SavedPost{})
	db.AutoMigrate(&domain.PostReports{})
	db.AutoMigrate(&domain.ViewStory{})
	
	return db, nil
}
