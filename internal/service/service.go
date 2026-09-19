package service

import (
	"time"

	"github.com/yukay/CRM/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type ServiceSupport struct {
	repo *repository.ReposirotySupport
}

func NewService(repo *repository.ReposirotySupport) *ServiceSupport {
	return &ServiceSupport{repo: repo}
}

func (s *ServiceSupport) RegistorUser(email, password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	createdAt := time.Now()

	return s.repo.CreateUser(email, string(bytes), createdAt)

}
