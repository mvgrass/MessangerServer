package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IUserService interface {
	healthHandler(*gin.Context)
	registerHandler(*gin.Context)
	loginHandler(*gin.Context)
}

type UserService struct {
	repo IUserRepository
}

func CreateUserService(repo IUserRepository) *UserService {
	return &UserService{repo: repo}
}

func (r *UserService) healthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (r *UserService) registerHandler(ctx *gin.Context) {
	var requestDto RegisterRequestDto
	if err := ctx.ShouldBindBodyWithJSON(&requestDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responseDto, err := r.repo.createUser(&requestDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, responseDto)
}

func (r *UserService) loginHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
