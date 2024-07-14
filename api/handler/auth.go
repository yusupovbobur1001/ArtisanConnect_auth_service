package handler

import (
	token "auth_service/api/token"
	pb "auth_service/genproto/auth"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Register User
// @Description Registers user
// @Tags Auth
// @ID register
// @Accept json
// @Produce json
// @Param user body auth.User true "User information to create it"
// @Success 201
// @Failure 500 {object} models.Error "Something went wrong in server"
// @Router /auth/register [post]
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
// @Failure 401 {object} models.Error "if Access token fails it will returns this"
// @Failure 500 {object} models.Error "Something went wrong in server"
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	user := pb.UserLogin{}

	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Logger.Error(err.Error())
		return
	}

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
// @Success 200 
// @Failure 401 {object} models.Error "if Access token fails it will returns this"
// @Failure 500 {object} models.Error "Something went wrong in server"
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	tkn := c.GetHeader("Authorization")

	_, err := token.ExtractClaims(tkn, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token",
		})
		h.Logger.Error("invalid token ", err.Error(), err)
		return
	}

	err = h.Auth.Logout(&pb.Tokens{RefreshToken: tkn})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, "SUCCESS")
}

// @Summary refreshes token
// @Description generates new access token gets token from header
// @Tags Auth
// @ID refresh
// @Produce json
// @Success 200 
// @Failure 500 {object} models.Error "Something went wrong in server"
// @Router /auth/refreshtoken [get]
func (h *Handler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")

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
		h.Logger.Error("err", err.Error(), err)
		return
	}

	accessToken := token.GenerateAccessToken(&claims)
	c.JSON(http.StatusOK, accessToken)
}
