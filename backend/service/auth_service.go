package service

import (
	"backend/repository"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService interface {
	Login(email, password string) (string, string, error)
}

type authService struct {
	repo repository.AuthRepository
}


func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Login(email, password string) (string, string, error) {
	// Retrieve user by email
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}

	//Check if the provided password matches the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Generate a new JWT token (access token) for the authenticated user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id": user.ID,
    "email":   user.Email,
    "exp":     time.Now().Add(30 * time.Minute).Unix(),
    "iat":     time.Now().Unix(),
})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {  
		return "", "", err
	}

	// We need a new refresh token for the user, so we return it to the caller
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte("your-refresh-secret-key"))
	if err != nil {
		return "", "", err
	}


	return tokenString, refreshTokenString, nil
}
