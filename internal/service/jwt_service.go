package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken = errors.New("token inválido ou expirado")
)

type JWTClaims struct {
	UserID  string `json:"user_id"`
	StoreID string `json:"store_id"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateAccessToken(userID, storeID, email string, ttl time.Duration) (string, error)
	ValidateAccessToken(tokenString string) (*JWTClaims, error)
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) bool
	HashToken(token string) string
}

type jwtService struct {
	secretKey []byte
}

func NewJWTService(secret string) JWTService {
	return &jwtService{
		secretKey: []byte(secret),
	}
}

func (j *jwtService) GenerateAccessToken(userID, storeID, email string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := &JWTClaims{
		UserID:  userID,
		StoreID: storeID,
		Email:   email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "cobertura085-api",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("falha ao assinar access token: %w", err)
	}

	return tokenString, nil
}

func (j *jwtService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", t.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (j *jwtService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("falha ao gerar hash da senha: %w", err)
	}
	return string(bytes), nil
}

func (j *jwtService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (j *jwtService) HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
