package token

import (
	"auth_service/config"
	pb "auth_service/genproto/auth"
	"auth_service/model"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GeneratorJWT(user *model.User) *pb.Tokens {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	claims := accessToken.Claims.(jwt.MapClaims)
	claims["user_name"] = user.UserName
	claims["full_name"] = user.FullName
	claims["bio"] = user.Bio
	claims["user_type"] = user.UserType
	claims["iat"] = time.Now().Unix()
	claims["ext"] = time.Now().Add(time.Minute).Unix()
	s := time.Now().Add(time.Minute).Unix()
	cfg := config.Load()

	access, err := accessToken.SignedString([]byte(cfg.REFRESH_SIGNING_KEY))
	if err != nil {
		log.Fatalf("Access token is not generated %v", err)
	}

	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["user_name"] = user.UserName
	refreshClaims["full_name"] = user.FullName
	refreshClaims["user_type"] = user.UserType
	refreshClaims["bio"] = user.Bio
	refreshClaims["iat"] = time.Now().Unix()
	refreshClaims["ext"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	refresh, err := refreshToken.SignedString([]byte(cfg.AUTH_SERVICE_PORT))
	if err != nil {
		log.Fatal("Refresh token is not generated %v", err)
	}

	return &pb.Tokens{
		AccessToken: access,
		RefreshToken: refresh,
		ExpiresIn: s,
	}
}
