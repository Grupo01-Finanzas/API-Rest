package controller

import (
	"fmt"
	"net/http"
	"strings"

	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthController handles authentication-related endpoints.
type AuthController struct {
	authService service.AuthService
}

// NewAuthController creates a new instance of AuthController.
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// RegisterAdmin godoc
// @Summary      Register Admin
// @Description  Registers a new admin user along with their establishment.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        registration  body      request.CreateAdminAndEstablishmentRequest  true  "Admin and establishment registration data"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Router       /register [post]
func (c *AuthController) RegisterAdmin(ctx *gin.Context) {
	var req request.CreateAdminAndEstablishmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if err := c.authService.RegisterAdmin(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Admin and establishment registered successfully"})
}

// Login godoc
// @Summary      Login
// @Description  Logs in a user with their email and password.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      request.LoginRequest  true  "User login credentials"
// @Success      200  {object}  response.AuthResponse
// @Failure      401  {object}  response.ErrorResponse
// @Router       /login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	authResponse, err := c.authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, authResponse)
}

// RefreshToken godoc
// @Summary      Refresh Token
// @Description  Refreshes the access token using a valid refresh token.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {refreshToken}"
// @Success      200  {object}  response.AuthResponse
// @Failure      401  {object}  response.ErrorResponse
// @Router       /refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Authorization header missing"})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Invalid Authorization header format"})
		return
	}

	refreshToken := parts[1]
	if refreshToken == "" {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Refresh token missing"})
		return
	}

	authResponse, err := c.authService.AttemptRefresh(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, authResponse)
}

// ResetPassword godoc
// @Summary      Reset Password
// @Description  Resets the password for the authenticated user.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        reset         body      request.ResetPasswordRequest  true  "Reset password request data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Router       /reset-password [post]
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req request.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Unauthorized"})
		return
	}

	claimsMap, ok := claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Internal server error: invalid claims"})
		return
	}

	userIDFloat, ok := claimsMap["user_id"].(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Internal server error: missing user ID"})
		return
	}

	userID := uint(userIDFloat)
	fmt.Println("User ID: ", userID) // Debug: Print userID

	err := c.authService.ResetPassword(&req, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
