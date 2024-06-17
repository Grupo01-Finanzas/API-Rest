package controller

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/service"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register godoc
// @Summary      Register a new user
// @Description  Registers a new user with the provided data
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      request.CreateUserRequest  true  "User data"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Router       /register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req request.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.authService.RegisterUser(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login godoc
// @Summary      Login
// @Description  Logs in with user credentials
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      request.LoginRequest  true  "User credentials"
// @Success      200  {object}  response.AuthResponse
// @Failure      401  {object}  response.ErrorResponse
// @Router      /login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authResponse, err := c.authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, authResponse)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Gets a new access token if the current access token has expired less than 5 minutes ago.
// @Description  If the access token is invalid or expired for more than 5 minutes, the user must log in again.
// @Description  The access token must be sent in the Authorization header as a Bearer token.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Success      200  {object}  response.AuthResponse
// @Failure      401  {object}  response.ErrorResponse
// @Router       /refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	accessToken := parts[1]
	if accessToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Access token missing"})
		return
	}

	authResponse, err := c.authService.AttemptRefresh(accessToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, authResponse)
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Resets the password for the authenticated user, requiring the current password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        reset         body      request.ResetPasswordRequest  true  "Reset password request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Router       /reset-password [post]
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req request.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claimsMap, ok := claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userIDFloat, ok := claimsMap["user_id"].(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userID := uint(userIDFloat)

	err := c.authService.ResetPassword(&req, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
