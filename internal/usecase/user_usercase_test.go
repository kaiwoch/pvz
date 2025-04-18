package usecase_test

import (
	"database/sql"
	"errors"
	"pvz/internal/storage/migrations/entity"
	"pvz/internal/usecase"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUsersStorage struct {
	mock.Mock
}

func (m *MockUsersStorage) GetUserByEmail(email string) (*entity.User, bool, error) {
	args := m.Called(email)
	return args.Get(0).(*entity.User), args.Bool(1), args.Error(2)
}

func (m *MockUsersStorage) CreateUser(email, password, role string) (*entity.User, error) {
	args := m.Called(email, password, role)
	return args.Get(0).(*entity.User), args.Error(1)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GenerateToken(userID uuid.UUID, role string) (string, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func TestUserUsecase_Login(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	role := "moderator"
	token := "test_token"

	tests := []struct {
		name          string
		email         string
		password      string
		mockUser      *entity.User
		mockUserErr   error
		mockToken     string
		mockTokenErr  error
		expected      string
		expectedError string
	}{
		{
			name:     "success",
			email:    email,
			password: password,
			mockUser: &entity.User{
				ID:       userID,
				Email:    email,
				Password: string(hashedPassword),
				Role:     role,
			},
			mockToken: token,
			expected:  token,
		},
		{
			name:          "user not found",
			email:         email,
			password:      password,
			mockUserErr:   sql.ErrNoRows,
			expectedError: "invalid credentials",
		},
		{
			name:     "wrong password",
			email:    email,
			password: "wrong_password",
			mockUser: &entity.User{
				ID:       userID,
				Email:    email,
				Password: string(hashedPassword),
				Role:     role,
			},
			expectedError: "invalid credentials",
		},
		{
			name:     "token generation error",
			email:    email,
			password: password,
			mockUser: &entity.User{
				ID:       userID,
				Email:    email,
				Password: string(hashedPassword),
				Role:     role,
			},
			mockTokenErr:  errors.New("token error"),
			expectedError: "token error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userStorage := new(MockUsersStorage)
			authService := new(MockAuthService)
			usecase := usecase.NewUserUsecase(userStorage, authService)

			userStorage.On("GetUserByEmail", tt.email).Return(tt.mockUser, false, tt.mockUserErr)

			if tt.mockUser != nil && tt.mockUserErr == nil && tt.expectedError != "invalid credentials" {
				authService.On("GenerateToken", tt.mockUser.ID, tt.mockUser.Role).Return(tt.mockToken, tt.mockTokenErr)
			}

			token, err := usecase.Login(tt.email, tt.password)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, token)
			}

			userStorage.AssertExpectations(t)
			authService.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_Register(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())
	email := "test@example.com"
	password := "password123"
	role := "user"

	tests := []struct {
		name          string
		email         string
		password      string
		role          string
		empty         bool
		mockGetUser   *entity.User
		mockGetErr    error
		mockCreateErr error
		expected      string
		expectedError string
	}{
		{
			name:        "success",
			email:       email,
			password:    password,
			role:        role,
			empty:       false,
			mockGetUser: nil,
			mockGetErr:  sql.ErrNoRows,
			expected:    userID.String(),
		},
		{
			name:          "user exists",
			email:         email,
			password:      password,
			role:          role,
			empty:         true,
			mockGetUser:   &entity.User{},
			expectedError: "user exist",
		},
		{
			name:          "get user error",
			email:         email,
			password:      password,
			role:          role,
			empty:         false,
			mockGetErr:    errors.New("db error"),
			expectedError: "error: db error",
		},
		{
			name:          "create user error",
			email:         email,
			password:      password,
			role:          role,
			empty:         false,
			mockGetErr:    sql.ErrNoRows,
			mockCreateErr: errors.New("create error"),
			expectedError: "create error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userStorage := new(MockUsersStorage)
			authService := new(MockAuthService)
			usecase := usecase.NewUserUsecase(userStorage, authService)

			userStorage.On("GetUserByEmail", tt.email).Return(tt.mockGetUser, tt.empty, tt.mockGetErr)

			if tt.mockGetErr == sql.ErrNoRows && tt.expectedError != "user exists" {
				userStorage.On("CreateUser", tt.email, mock.Anything, tt.role).
					Return(&entity.User{ID: userID}, tt.mockCreateErr)
			}

			result, err := usecase.Register(tt.email, tt.password, tt.role)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			userStorage.AssertExpectations(t)
		})
	}
}
