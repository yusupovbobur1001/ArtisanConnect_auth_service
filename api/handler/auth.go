package handler

import (
	"auth_service/api/email"
	token "auth_service/api/token"
	pb "auth_service/genproto/auth"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
)

// @Summary Register User
// @Description Registers user
// @Tags Auth
// @ID register
// @Accept json
// @Produce json
// @Param user body auth.User true "User information to create it"
// @Success 201  string   "SUCCESS"
// @Failure 500  string    auth.Message
// @Router /register [post]
func (h *Handler) Register(c *gin.Context) {
	user := &pb.User{}
	if err := c.BindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Logger.Error(err.Error())
		return
	}

	err := h.Auth.Register(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusCreated, "SUCCESS")
}

// @Summary Login user
// @Description checks the user and returns tokens
// @Tags Auth
// @ID login
// @Accept json
// @Produce json
// @Param user body auth.UserLogin true "User Information to log in"
// @Success 200 {object} auth.Tokens  "Returns access and refresh tokens"
// @Failure 401 {object} string "if Access token fails it will returns this"
// @Failure 500 {object} string "Something went wrong in server"
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	user := pb.UserLogin{}

	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Logger.Error(err.Error())
		return
	}

	fmt.Println(user.Email, user.Password)

	resp, err := h.Auth.Login(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user does not exist"})
			h.Logger.Error(err.Error())
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			h.Logger.Error(err.Error())
			return
		}
	}

	auth := token.GeneratorJWT(resp)

	c.JSON(http.StatusOK, auth)
}

// @Summary log outs user
// @Description removes refresh token gets token from header
// @Tags Auth
// @ID logout
// @Accept json
// @Produce json
// @Success 200  string  "SUCCESS"
// @Failure 401 {object} string "if Access token fails it will returns this"
// @Failure 500 {object} string "Something went wrong in server"
// @Router /logout [post]
func (h *Handler) Logout(c *gin.Context) {
	tkn := c.GetHeader("Authorization")

	_, err := token.ExtractClaims(tkn, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token",
		})
		h.Logger.Error("invalid token..... ", err.Error(), err)
		return
	}

	err = h.Auth.Logout(&pb.Tokens{RefreshToken: tkn})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		h.Logger.Error("asfd.....", "error: ", err.Error())	
		return
	}

	c.JSON(http.StatusOK, "SUCCESS")
}

// @Summary refreshes token
// @Description generates new access token gets token from header
// @Tags Auth
// @ID refresh
// @Produce json
// @Param token body auth.Refreshtoken true "Token"
// @Success 200  {object} string "if Access token fails it will returns this"
// @Failure 500 {object} string "Something went wrong in server"
// @Router /refreshtoken [get]
func (h *Handler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	fmt.Println("salom refreshToken")
	claims, err := token.ExtractClaims(refreshToken, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token",
		})
		h.Logger.Error(err.Error())
		return
	}

	exist, err := h.Auth.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		h.Logger.Error(err.Error())
		return
	}

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token doesnt exists",
		})
		h.Logger.Error("err token yoq", "err: ", err)
		return
	}

	accessToken := token.GenerateAccessToken(&claims)
	c.JSON(http.StatusOK, accessToken)
}

// Passwordrecovery godoc
// @Summary Recover password
// @Description Send password recovery email
// @Tags auth
// @Accept  json
// @Produce  json
// @Param  body  body  auth.RestoreProfile  true  "Password Recovery Request"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /passwordrecovery [post]
func (h *Handler) Passwordrecovery(c *gin.Context) {
	req := pb.RestoreProfile{}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	rand.Seed(uint64(time.Now().Unix()))
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	ctx := context.Background()

	err = h.Redis.Set(ctx, req.Email, code, time.Minute*8).Err()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	email.SendCode(req.Email, code)
	c.JSON(http.StatusAccepted, gin.H{"message": "Password recovery email sent"})
}
