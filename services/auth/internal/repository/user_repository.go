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
	CreateUser(*model.User) error
	GetUserByEmail(string) (*model.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) CreateUser(user *model.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("Email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
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

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		fmt.Println("Migration error:", err)
	}
	return &UserRepository{db: db}
}
