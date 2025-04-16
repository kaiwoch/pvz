package usecase

import (
	"database/sql"
	"fmt"
	"pvz/internal/storage"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userStorage *storage.UsersPostgresStorage
	authService *AuthService
}

func NewUserUsecase(userStorage *storage.UsersPostgresStorage, authService *AuthService) *UserUsecase {
	return &UserUsecase{userStorage: userStorage, authService: authService}
}

func (u *UserUsecase) Login(email, password string) (string, error) {
	user, err := u.userStorage.GetUserByEmail(email)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	return u.authService.GenerateToken(user.ID, user.Role)
}

func (u *UserUsecase) Register(email, password, role string) (string, error) {
	user, err := u.userStorage.GetUserByEmail(email)
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("user exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to check user existence: %w", err)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user, err = u.userStorage.CreateUser(email, string(hashedPassword), role)
	if err != nil {
		return "", err
	}
	return user.ID.String(), nil
}
