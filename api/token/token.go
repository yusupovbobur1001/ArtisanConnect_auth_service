package token

import (
	"auth_service/config"
	pb "auth_service/genproto/auth"
	"auth_service/model"
	"fmt"
	"log"
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
	fmt.Println(cfg.SIGNING_KEY)
	access, err := accessToken.SignedString([]byte(cfg.SIGNING_KEY))
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
	fmt.Println(cfg.REFRESH_SIGNING_KEY)
	refresh, err := refreshToken.SignedString([]byte(cfg.REFRESH_SIGNING_KEY))
	if err != nil {
		log.Fatalf("Refresh token is not generated %v", err)
	}	
	fmt.Println(access)
	fmt.Println("---------------------------")
	fmt.Println(refresh)
	return &pb.Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    s,
	}
}

func GenerateAccessToken(user *jwt.MapClaims) *string {

	accessToken := jwt.New(jwt.SigningMethodHS256)

	claims := accessToken.Claims.(jwt.MapClaims)
	claims["user_name"] = (*user)["user_name"]
	claims["full_name"] = (*user)["full_name"]
	claims["bio"] = (*user)["bio"]
	claims["user_type"] = (*user)["user_type"]
	claims["iat"] = time.Now().Unix()
	claims["ext"] = time.Now().Add(time.Minute).Unix()

	cfg := config.Load()

	access, err := accessToken.SignedString([]byte(cfg.SIGNING_KEY))
	if err != nil {
		log.Fatalf("Access token is not generated %v", err)
	}
	fmt.Println(access, "++++++++++")
	return &access
}



func ExtractClaims(tokenStr string, isRefresh bool) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		if isRefresh {
			return []byte(config.Load().REFRESH_SIGNING_KEY), nil
		}
		return []byte(config.Load().SIGNING_KEY), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
