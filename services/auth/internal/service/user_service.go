package service

import (
	"MessangerServer/services/auth/internal/model"
	"MessangerServer/services/auth/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IUserService interface {
	HealthHandler(*gin.Context)
	RegisterHandler(*gin.Context)
	LoginHandler(*gin.Context)
	LogoutHandler(*gin.Context)
	RefreshHandler(*gin.Context)
	GetMyselfHandler(*gin.Context)
}

type UserService struct {
	repo repository.IUserRepository
}

func CreateUserService(repo repository.IUserRepository) *UserService {
	return &UserService{repo: repo}
}

func (r *UserService) HealthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (r *UserService) RegisterHandler(ctx *gin.Context) {
	var requestDto RegisterRequestDto
	if err := ctx.ShouldBindBodyWithJSON(&requestDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if user already exist
	err := r.repo.CreateUser(&model.User{Name: requestDto.Name, Email: requestDto.Email})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, RegisterResponseDto{AccessToken: "placeholder", RefreshToken: "placeholder"})
}

func (r *UserService) LoginHandler(ctx *gin.Context) {
	// check db and password
	// create access token create refresh token

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (r *UserService) LogoutHandler(ctx *gin.Context) {
	// remove refresh token from db
	// optional create and add to blacklist current access token until exparation
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (r *UserService) RefreshHandler(ctx *gin.Context) {
	// return new access token
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (r *UserService) GetMyselfHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
