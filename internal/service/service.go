package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yukay/CRM/internal/models"
	"github.com/yukay/CRM/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type ServiceSupport struct {
	repo      *repository.ReposirotySupport
	secretKey []byte
}

func NewService(repo *repository.ReposirotySupport, secretkey []byte) *ServiceSupport {
	return &ServiceSupport{repo: repo, secretKey: secretkey}
}

func (s *ServiceSupport) RegistorUser(email, password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	createdAt := time.Now()

	return s.repo.CreateUser(email, string(bytes), createdAt)

}

func (s *ServiceSupport) CheckLoginUser(email, password string) (int, error) {

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return 0, err
	}
	return user.Id, nil

}

func (s *ServiceSupport) GenerateJWT(userId int) (string, error) {

	claims := models.Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "CRM",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (s *ServiceSupport) LoginUser(email, password string) (string, error) {

	userID, err := s.CheckLoginUser(email, password)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := s.GenerateJWT(userID)
	if err != nil {
		return "", err
	}

	return token, nil

}
