package service

import (
	"MessangerServer/services/auth/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IUserService interface {
	HealthHandler(*gin.Context)
	RegisterHandler(*gin.Context)
	LoginHandler(*gin.Context)
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
	var requestDto repository.RegisterRequestDto
	if err := ctx.ShouldBindBodyWithJSON(&requestDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responseDto, err := r.repo.CreateUser(&requestDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, responseDto)
}

func (r *UserService) LoginHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
