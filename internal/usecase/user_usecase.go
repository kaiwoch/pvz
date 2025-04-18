package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"pvz/internal/storage"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Login(email, password string) (string, error)
	Register(email, password, role string) (string, error)
}

type UserUsecaseImpl struct {
	userStorage storage.UsersPostgresStorage
	authService AuthUsecase
}

func NewUserUsecase(userStorage storage.UsersPostgresStorage, authService AuthUsecase) *UserUsecaseImpl {
	return &UserUsecaseImpl{userStorage: userStorage, authService: authService}
}

func (u *UserUsecaseImpl) Login(email, password string) (string, error) {
	user, _, err := u.userStorage.GetUserByEmail(email)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	return u.authService.GenerateToken(user.ID, user.Role)
}

func (u *UserUsecaseImpl) Register(email, password, role string) (string, error) {
	user, ok, err := u.userStorage.GetUserByEmail(email)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("error: %w", err)
	}

	if ok {
		return "", errors.New("user exist")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user, err = u.userStorage.CreateUser(email, string(hashedPassword), role)
	if err != nil {
		return "", err
	}
	return user.ID.String(), nil
}
