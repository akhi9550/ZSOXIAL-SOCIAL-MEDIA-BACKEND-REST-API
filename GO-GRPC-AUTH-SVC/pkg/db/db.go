package db

import (
    "fmt"

    "github.com/akhi9550/auth-svc/pkg/config"
    "github.com/akhi9550/auth-svc/pkg/domain"
    "github.com/akhi9550/auth-svc/pkg/helper"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
    psqlInfo := fmt.Sprintf("host=%s user=%s dbname=postgres port=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPort, cfg.DBPassword)
    db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
    if err != nil {
        return nil, err 
    }

    rows, err := db.Raw("SELECT 1 FROM pg_database WHERE datname = ?", cfg.DBName).Rows()
    if err != nil {
        return nil, err 
    }
    defer rows.Close()

    if !rows.Next() {
        query := fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
        if err := db.Exec(query).Error; err != nil {
            return nil, err
        }
        fmt.Printf("Database %s created successfully.\n", cfg.DBName)
    }
    psqlInfo = fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
    db, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    db.AutoMigrate(
        &domain.User{},
        &domain.UserReports{},
        &domain.FollowingRequests{},
        &domain.Followings{},
        &domain.Followers{},
        &domain.Groups{},
        &domain.GroupMembers{},
    )

    CreateDefaultAdmin(db)

    return db, nil
}

func CreateDefaultAdmin(db *gorm.DB) {
    var count int64
    db.Model(&domain.User{}).Count(&count)
    if count == 0 {
        password := "admin@123"
        hashPassword, err := helper.PasswordHash(password)
        if err != nil {
            return
        }
        admin := domain.User{
            ID:        1,
            Firstname: "Zsoxial",
            Lastname:  "Admin",
            Username:  "admin",
            Dob:       "10-10-2000",
            Gender:    "male",
            Phone:     "+919061757507",
            Email:     "admin@zsoxial.com",
            Password:  hashPassword,
            Bio:       "",
            Imageurl:  "https://img.freepik.com/free-psd/3d-illustration-human-avatar-profile_23-2150671142.jpg",
            Blocked:   false,
            Isadmin:   true,
        }
        db.Create(&admin)
    }
}
