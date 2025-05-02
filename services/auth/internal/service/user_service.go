package service

import (
	"MessangerServer/services/auth/internal/config"
	"MessangerServer/services/auth/internal/model"
	"MessangerServer/services/auth/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	repo   repository.IUserRepository
	pepper string
}

func CreateUserService(repo repository.IUserRepository, cfg *config.Config) *UserService {
	return &UserService{repo: repo, pepper: cfg.App.Pepper}
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

	if err := validator.New().Struct(requestDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation of input data failed"})
		return
	}

	//looks like need to sepate logic from transport
	if _, err := r.repo.GetUserByEmail(requestDto.Email); err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "This email already used!"})
		return
	}

	user := model.User{Uuid: uuid.New().String(), Name: requestDto.Name, Email: requestDto.Email}
	pass, err := bcrypt.GenerateFromPassword([]byte(requestDto.Password+r.pepper), 14)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password is too long"})
		return
	}
	user.Password = string(pass)

	err = r.repo.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, UserInfoResponseDto{Id: user.Uuid, Name: user.Name, Email: user.Email})
}

func (r *UserService) LoginHandler(ctx *gin.Context) {
	var requestDto LoginRequestDto
	if err := ctx.ShouldBindBodyWithJSON(&requestDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userInfo, err := r.repo.GetUserByEmail(requestDto.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(requestDto.Password+r.pepper)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong email or password"})
		return
	}

	ctx.JSON(http.StatusOK, JwtTokenRespnseDto{})
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
