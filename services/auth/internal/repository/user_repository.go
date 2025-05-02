package repository

import (
	"MessangerServer/services/auth/internal/config"
	"MessangerServer/services/auth/internal/model"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(*RegisterRequestDto) (*RegisterResponseDto, error)
}

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) CreateUser(request *RegisterRequestDto) (*RegisterResponseDto, error) {
	user := model.User{
		Name:  request.Name,
		Email: request.Email,
	}

	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &RegisterResponseDto{AccessToken: "AccessTokenPlaceholder", RefreshToken: "RefreshTokenPlaceHolder"}, nil
}

func InitStorage(cfg *config.Config) *UserRepository {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Db.Addr,
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.DbName,
		cfg.Db.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to PostgreSQL!")

	db.AutoMigrate(&model.User{})

	return &UserRepository{db: db}
}
