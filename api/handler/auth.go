package handler

import (
	"auth-service/api/tokens"
	"auth-service/models"
	"auth-service/storage/redis"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary Registers user
// @Description Registers a new user
// @Tags auth
// @Param user body models.RegisterRequest true "User data"
// @Success 200 {object} models.RegisterResponse
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /register [post]
func (h *Handler) Register(c *gin.Context) {
	h.Logger.Info("Register handler is invoked")

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, h, err, "invalid data", http.StatusBadRequest)
		return
	}

	passByte, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		handleError(c, h, err, "error hashing password", http.StatusInternalServerError)
		return
	}
	req.Password = string(passByte)

	ctx, cancel := context.WithTimeout(context.Background(), h.ContextTimeout)
	defer cancel()

	resp, err := h.Storage.User().Add(ctx, &req)
	if err != nil {
		handleError(c, h, err, "error registering user", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Register handler is completed successfully")
	c.JSON(http.StatusOK, resp)
}

// Login godoc
// @Summary Login user
// @Description Logs user in
// @Tags auth
// @Param user body models.LoginReq true "User data"
// @Success 200 {object} models.Tokens
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	h.Logger.Info("Login handler is invoked")

	var req models.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, h, err, "invalid data", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.ContextTimeout)
	defer cancel()

	user, err := h.Storage.User().GetDetails(ctx, req.Email)
	if err != nil {
		handleError(c, h, err, "error getting user details", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		handleError(c, h, err, "invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, err := tokens.GenerateAccessToken(h.Config, user.Id, user.Role)
	if err != nil {
		handleError(c, h, err, "error generating access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := tokens.GenerateRefreshToken(h.Config, user.Id)
	if err != nil {
		handleError(c, h, err, "error generating refresh token", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Login handler is completed successfully")
	c.JSON(http.StatusOK, models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Logout godoc
// @Summary Logouts user
// @Description Logouts user
// @Tags auth
// @Param email query string true "User email"
// @Success 200 {string} string "User logged out successfully"
// @Failure 400 {object} string "Invalid email"
// @Failure 500 {object} string "Server error while processing request"
// @Router /logout [post]
func (h *Handler) Logout(c *gin.Context) {
	h.Logger.Info("Logout handler is invoked")

	email := c.Query("email")
	if email == "" {
		handleError(c, h, errors.New("no email provided"), "invalid email", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.ContextTimeout)
	defer cancel()

	user, err := h.Storage.User().GetDetails(ctx, email)
	if err != nil {
		handleError(c, h, err, "error getting user details", http.StatusInternalServerError)
		return
	}

	err = redis.DeleteToken(h.Config, ctx, user.Id)
	if err != nil {
		handleError(c, h, err, "error logging out", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Logout handler is completed successfully")
	c.JSON(http.StatusOK, "User logged out successfully")
}

// RefreshToken godoc
// @Summary Refreshes token
// @Description Refreshes refresh token
// @Tags auth
// @Param data body models.RefreshToken true "Refresh token"
// @Success 200 {object} models.Tokens
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	h.Logger.Info("RefreshToken handler is invoked")

	var req models.RefreshToken
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, h, err, "invalid data", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.ContextTimeout)
	defer cancel()

	valid, err := tokens.ValidateRefreshToken(h.Config, req.RefreshToken)
	if !valid || err != nil {
		handleError(c, h, err, "error validating refresh token", http.StatusInternalServerError)
		return
	}

	id, err := tokens.ExtractRefreshUserID(h.Config, req.RefreshToken)
	if err != nil {
		handleError(c, h, err, "error extracting user id from refresh token", http.StatusInternalServerError)
		return
	}

	role, err := h.Storage.User().GetRole(ctx, id)
	if err != nil {
		handleError(c, h, err, "error getting user role", http.StatusInternalServerError)
		return
	}

	accessToken, err := tokens.GenerateAccessToken(h.Config, id, role)
	if err != nil {
		handleError(c, h, err, "error generating access token", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("RefreshToken handler is completed successfully")
	c.JSON(http.StatusOK, models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
	})
}

// ValidateToken godoc
// @Summary Validates token
// @Description Validates access token
// @Tags auth
// @Param data body models.AccessToken true "Access token"
// @Success 200 {string} string "Access token is valid"
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /validate [post]
func (h *Handler) ValidateToken(c *gin.Context) {
	h.Logger.Info("ValidateToken handler is invoked")

	var req models.AccessToken
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, h, err, "invalid data", http.StatusBadRequest)
		return
	}

	valid, err := tokens.ValidateAccessToken(h.Config, req.AccessToken)
	if !valid || err != nil {
		handleError(c, h, err, "error validating access token", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("ValidateToken handler is completed successfully")
	c.JSON(http.StatusOK, "Access token is valid")
}
