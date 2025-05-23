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
	GetUserByUserId(string) (*model.User, error)
	StoreRefreshToken(string, string) error
	GetRefreshToken(string) (string, error)
	RevokeTokens(string, string) error
	CacheUserByToken(string, string, time.Duration) error
	IsTokenRevoked(string) (int64, error)
	GetUserIdByAccessTokenCache(string) (string, error)
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

func (r *AuthRepository) GetUserByUserId(userId string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("Uuid = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("Email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) StoreRefreshToken(userId, refreshToken string) error {
	ctx := context.Background()

	return r.redisClient.Set(ctx, fmt.Sprintf("auth:refresh:%s", userId), refreshToken, 7*24*time.Hour).Err()
}

func (r *AuthRepository) GetRefreshToken(userId string) (string, error) {
	ctx := context.Background()
	return r.redisClient.Get(ctx, fmt.Sprintf("auth:refresh:%s", userId)).Result()
}

func (r *AuthRepository) RevokeTokens(userId, accessToken string) error {
	ctx := context.Background()
	// TODO: Move time to config for centrilized usage
	return r.redisClient.Set(ctx, fmt.Sprintf("auth:revoked:%s", accessToken), "1", 15*time.Minute).Err()
}

func (r *AuthRepository) CacheUserByToken(accessToken, userId string, exparation time.Duration) error {
	ctx := context.Background()
	return r.redisClient.Set(ctx, fmt.Sprintf("auth:accessCache:%s", accessToken), userId, exparation).Err()
}

func (r *AuthRepository) IsTokenRevoked(accessToken string) (int64, error) {
	ctx := context.Background()
	return r.redisClient.Exists(ctx, fmt.Sprintf("auth:revoked:%s", accessToken)).Result()
}

func (r *AuthRepository) GetUserIdByAccessTokenCache(accessToken string) (string, error) {
	ctx := context.Background()
	return r.redisClient.Get(ctx, fmt.Sprintf("auth:accessCache:%s", accessToken)).Result()
}

func InitStorage(cfg *config.Config) *AuthRepository {
	// TODO: think deeper about work with ctx
	var repo AuthRepository
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		opt, err := redis.ParseURL(cfg.RedisUrl)
		if err != nil {
			fmt.Println(err)
			return err
		}
		repo.redisClient = redis.NewClient(opt)
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err = repo.redisClient.Ping(pingCtx).Err()
		if err != nil {
			fmt.Println("Can't connect to redis instance")
			return err
		}

		fmt.Println("Successfully connceted to redis!")

		return nil
	})

	group.Go(func() error {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.Db.Addr,
			cfg.Db.User,
			cfg.Db.Password,
			cfg.Db.DbName,
			cfg.Db.Port)

		// TODO: split on postgress connection
		// and create if not exist
		var err error
		repo.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println(err)
			return err
		}
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		sqlDb, err := repo.db.DB()
		if err != nil {
			return err
		}
		if err = sqlDb.PingContext(pingCtx); err != nil {
			return err
		}
		fmt.Println("Connected to PostgreSQL!")

		err = repo.db.AutoMigrate(&model.User{})
		if err != nil {
			fmt.Println("Migration error:", err)
			return err
		}

		return nil
	})

	if err := group.Wait(); err != nil {
		if repo.redisClient != nil {
			repo.redisClient.Close()
		}
		if repo.db != nil {
			sqlDb, err := repo.db.DB()
			if err != nil {
				sqlDb.Close()
			}
		}
		log.Fatal()
	}

	return &repo
}
