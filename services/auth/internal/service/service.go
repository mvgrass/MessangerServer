package service

import (
	"MessangerServerAuth/internal/config"
	"MessangerServerAuth/internal/model"
	"MessangerServerAuth/internal/repository"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	HealthHandler(*gin.Context)
	RegisterHandler(*gin.Context)
	LoginHandler(*gin.Context)
	LogoutHandler(*gin.Context)
	RefreshHandler(*gin.Context)
	GetMyselfHandler(*gin.Context)
	Validate(*gin.Context)
}

type AuthService struct {
	repo             repository.IAuthRepository
	pepper           string
	accessJwtSecret  string
	refreshJwtSecret string
	bcryptCost       int
}

func CreateAuthService(repo repository.IAuthRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		repo:             repo,
		pepper:           cfg.App.Pepper,
		bcryptCost:       cfg.App.Cost,
		accessJwtSecret:  cfg.JwtTokenSecret,
		refreshJwtSecret: cfg.JwtTokenSecret}
}

func (r *AuthService) HealthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (r *AuthService) RegisterHandler(ctx *gin.Context) {
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
	pass, err := bcrypt.GenerateFromPassword([]byte(requestDto.Password+r.pepper), r.bcryptCost)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password is too long"})
		return
	}
	user.Password = string(pass)

	err = r.repo.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, UserInfoResponseDto{Id: user.Uuid, Name: user.Name, Email: user.Email})
}

func (r *AuthService) LoginHandler(ctx *gin.Context) {
	var requestDto LoginRequestDto
	if err := ctx.ShouldBindBodyWithJSON(&requestDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userInfo, err := r.repo.GetUserByEmail(requestDto.Email)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(requestDto.Password+r.pepper)); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong email or password"})
		return
	}

	accessToken, refreshToken, err := GenerateToken(userInfo.Uuid, r.accessJwtSecret, r.refreshJwtSecret)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Can't generate JWT"})
		return
	}

	r.repo.StoreRefreshToken(userInfo.Uuid, refreshToken)

	ctx.JSON(http.StatusOK, JwtTokenRespnseDto{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresIn:  int64(15 * time.Minute),
		RefreshExpiresIn: int64(7 * 24 * time.Hour),
	})
}

func (r *AuthService) LogoutHandler(ctx *gin.Context) {
	accessToken := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
	if accessToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userClaims, err := ParseAccessTokenWithoutExparation(accessToken, r.accessJwtSecret)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing your token!"})
		return
	}

	err = r.repo.RevokeTokens(userClaims.UserId, accessToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request!"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "Loged out!",
	})
}

func (r *AuthService) RefreshHandler(ctx *gin.Context) {
	var reqDto RefreshRequestDto
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	oldAccessToken := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
	oldRefreshToken := reqDto.RefreshToken

	refreshClaims, err := ParseRefreshToken(oldRefreshToken, r.refreshJwtSecret)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong or expired refresh token!"})
		return
	}

	storedToken, err := r.repo.GetRefreshToken(refreshClaims.Subject)
	if storedToken != oldRefreshToken || err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong or expired refresh token!"})
		return
	}

	if err = r.repo.RevokeTokens(refreshClaims.Subject, oldAccessToken); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong or expired refresh token!"})
		return
	}

	accessToken, refreshToken, err := GenerateToken(refreshClaims.Subject, r.accessJwtSecret, r.refreshJwtSecret)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Can't generate JWT"})
		return
	}

	r.repo.StoreRefreshToken(refreshClaims.Subject, refreshToken)

	ctx.JSON(http.StatusOK, JwtTokenRespnseDto{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresIn:  int64(15 * time.Minute),
		RefreshExpiresIn: int64(7 * 24 * time.Hour),
	})
}

func (r *AuthService) GetMyselfHandler(ctx *gin.Context) {
	userId := ctx.GetHeader("X-User-ID")
	user, err := r.repo.GetUserByUserId(userId)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request!"})
		return
	}

	ctx.JSON(http.StatusOK, UserInfoResponseDto{Id: user.Uuid, Name: user.Name, Email: user.Email})
}

func (r *AuthService) Validate(ctx *gin.Context) {
	accessToken := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
	if accessToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if revoked, _ := r.repo.IsTokenRevoked(accessToken); revoked != 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if userId, err := r.repo.GetUserIdByAccessTokenCache(accessToken); err == nil {
		ctx.Header("X-User-ID", userId)
		ctx.Status(http.StatusOK)
		return
	}

	userClaims, err := ParseAccessToken(accessToken, r.accessJwtSecret)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	r.repo.CacheUserByToken(accessToken, userClaims.UserId, time.Until(userClaims.ExpiresAt.Time))

	ctx.Header("X-User-ID", userClaims.UserId)
	ctx.Status(http.StatusOK)
}
