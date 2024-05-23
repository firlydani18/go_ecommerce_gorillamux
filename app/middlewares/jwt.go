package middlewares

import (
	"context"
	"fmt"
	"go-ecommerce/app/config"
	"go-ecommerce/model"
	"go-ecommerce/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "user_id"

func CreateJWT(secret []byte, UserId int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpire)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    strconv.Itoa(UserId),
		"expired_at": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store model.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get token from user req
		tokenString := getTokenFromRequest(r)

		// validate jwt
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		// if is we need to fetch the user id from db (id from token)
		claims := token.Claims.(jwt.MapClaims)
		str := claims["user_id"].(string)
		userId, _ := strconv.Atoi(str)
		u, err := store.GetUserById(userId)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
		}
		permissionDenied(w)
		return
		// set context "userid" to the user id
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token == "" {
		return token
	}
	return ""
}

func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIdFromCtx(ctx context.Context) int {
	userId, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userId
}
