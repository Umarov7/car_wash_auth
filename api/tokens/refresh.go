package tokens

import (
	"auth-service/config"
	"auth-service/storage/redis"
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

func GenerateRefreshToken(cfg *config.Config, userID string) (string, error) {
	token := *jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	newToken, err := token.SignedString([]byte(cfg.REFRESH_TOKEN_KEY))
	if err != nil {
		log.Println(err)
		return "", errors.Wrap(err, "failed to generate refresh token")
	}

	err = redis.StoreToken(cfg, context.Background(), userID, newToken)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func ValidateRefreshToken(cfg *config.Config, tokenStr string) (bool, error) {
	_, err := ExtractRefreshClaims(cfg, tokenStr)
	if err != nil {
		return false, errors.Wrap(err, "validation failed")
	}

	id, err := ExtractRefreshUserID(cfg, tokenStr)
	if err != nil {
		return false, errors.Wrap(err, "validation failed")
	}

	token, err := redis.GetToken(cfg, context.Background(), id)
	if err != nil {
		return false, errors.Wrap(err, "validation failed")
	}

	if token != tokenStr {
		return false, errors.New("invalid refresh token")
	}

	return true, nil
}

func ExtractRefreshClaims(cfg *config.Config, tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.REFRESH_TOKEN_KEY), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse refresh token")
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func ExtractRefreshUserID(cfg *config.Config, tokenStr string) (string, error) {
	claims, err := ExtractRefreshClaims(cfg, tokenStr)
	if err != nil {
		return "", errors.Wrap(err, "failed to extract refresh token's claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid token claims: user id not found")
	}

	return userID, nil
}
