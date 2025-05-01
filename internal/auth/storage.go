package auth

import (
	authCfg "MessangerServer/internal/config/auth"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100;not null"`
	Email string `gorm:"size:255;uniqueIndex;not null"`
}

type IUserRepository interface {
	createUser(*RegisterRequestDto) (*RegisterResponseDto, error)
}

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) createUser(request *RegisterRequestDto) (*RegisterResponseDto, error) {
	user := User{
		Name:  request.Name,
		Email: request.Email,
	}

	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &RegisterResponseDto{AccessToken: "AccessTokenPlaceholder", RefreshToken: "RefreshTokenPlaceHolder"}, nil
}

func InitStorage(cfg *authCfg.Config) *UserRepository {

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

	db.AutoMigrate(&User{})

	return &UserRepository{db: db}
}
