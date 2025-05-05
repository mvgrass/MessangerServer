package repository

import (
	"MessangerServerAuth/internal/config"
	"MessangerServerAuth/internal/model"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	CreateUser(*model.User) error
	GetUserByEmail(string) (*model.User, error)
	//StoreRefreshToken(*model.User, string)
	//DeleteRefreshToken(string)
}

type AuthRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func (r *AuthRepository) CreateUser(user *model.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("Email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func InitStorage(cfg *config.Config) *AuthRepository {
	var redisClient *redis.Client
	var db *gorm.DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		opt, err := redis.ParseURL(cfg.RedisUrl)
		if err != nil {
			fmt.Println(err)
			return err
		}
		redisClient := redis.NewClient(opt)
		ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = redisClient.Ping(ctx1).Err()
		if err != nil {
			fmt.Println("Can't connect to redis instance")
			return err
		}

		fmt.Println("Successfully connceted to redis!")

		group.Go(func() error {
			<-ctx.Done()
			sqlDb, err := db.DB()
			if err != nil {
				return sqlDb.Close()
			}
			return err
		})

		return nil
	})

	group.Go(func() error {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.Db.Addr,
			cfg.Db.User,
			cfg.Db.Password,
			cfg.Db.DbName,
			cfg.Db.Port)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("Connected to PostgreSQL!")

		err = db.AutoMigrate(&model.User{})
		if err != nil {
			fmt.Println("Migration error:", err)
			return err
		}

		group.Go(func() error {
			<-ctx.Done()
			return redisClient.Close()
		})

		return nil
	})

	if err := group.Wait(); err != nil {
		if redisClient != nil {
			redisClient.Close()
		}
		if db != nil {
			sqlDb, err := db.DB()
			if err != nil {
				sqlDb.Close()
			}
		}
		log.Fatal()
	}

	return &AuthRepository{redisClient: redisClient, db: db}
}
