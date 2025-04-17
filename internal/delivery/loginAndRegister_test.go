package delivery_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pvz/internal/delivery"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) Register(email, password, role string) (string, error) {
	args := m.Called(email, password, role)
	return args.String(0), args.Error(1)
}

func TestDummyLoginHandler_DummyLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name         string
		requestBody  any
		mock         func(*MockUserUsecase)
		expectedCode int
	}{
		{
			name: "succes moderator login",
			requestBody: map[string]string{
				"role": "moderator",
			},
			mock: func(muu *MockUserUsecase) {
				muu.On("Login", "dummy_moderator@test.com", "supersecretpassword").Return("moderator_token", nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "succes employee login",
			requestBody: map[string]string{
				"role": "employee",
			},
			mock: func(muu *MockUserUsecase) {
				muu.On("Login", "dummy_employee@test.com", "supersecretpassword").Return("employee_token", nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid role",
			requestBody: map[string]string{
				"role": "invalid",
			},
			mock:         func(m *MockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing role",
			requestBody:  map[string]string{},
			mock:         func(m *MockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "moderator login error",
			requestBody: map[string]string{
				"role": "moderator",
			},
			mock: func(m *MockUserUsecase) {
				m.On("Login", "dummy_moderator@test.com", "supersecretpassword").
					Return("", errors.New("login failed"))
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "employee login error",
			requestBody: map[string]string{
				"role": "employee",
			},
			mock: func(m *MockUserUsecase) {
				m.On("Login", "dummy_employee@test.com", "supersecretpassword").
					Return("", errors.New("login failed"))
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "missing body",
			requestBody:  nil,
			mock:         func(m *MockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid json",
			requestBody:  "json",
			mock:         func(m *MockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &MockUserUsecase{}
			tt.mock(mockUsecase)

			handler := delivery.NewDummyLoginHandler(mockUsecase)

			router := gin.Default()
			router.POST("/dummyLogin", handler.DummyLogin)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "token")
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name         string
		requestBody  any
		mock         func(*MockUserUsecase)
		expectedCode int
	}{
		{
			name: "successful registration",
			requestBody: map[string]string{
				"email":    "kanzartem11@mail.ru",
				"password": "q1w2e3",
				"role":     "moderator",
			},
			mock: func(muu *MockUserUsecase) {
				muu.On("Register", "kanzartem11@mail.ru", "q1w2e3", "moderator").Return("user_id", nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "user already exists",
			requestBody: map[string]string{
				"email":    "kanzartem11@mail.ru",
				"password": "q1w2e3",
				"role":     "moderator",
			},
			mock: func(muu *MockUserUsecase) {
				muu.On("Register", "kanzartem11@mail.ru", "q1w2e3", "moderator").Return("", errors.New("user already exists"))
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing body",
			requestBody:  nil,
			mock:         func(muu *MockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "empty email",
			requestBody: map[string]string{
				"email":    "",
				"password": "q1w2e3",
				"role":     "moderator",
			},
			mock:         func(muu *MockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "empty password",
			requestBody: map[string]string{
				"email":    "kanzartem11@mail.ru",
				"password": "",
				"role":     "moderator",
			},
			mock:         func(muu *MockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "empty role",
			requestBody: map[string]string{
				"email":    "kanzartem11@mail.ru",
				"password": "q1w2e3",
				"role":     "",
			},
			mock:         func(muu *MockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid role",
			requestBody: map[string]string{
				"email":    "kanzartem11@mail.ru",
				"password": "q1w2e3",
				"role":     "user",
			},
			mock: func(muu *MockUserUsecase) {
				muu.On("Register", "kanzartem11@mail.ru", "q1w2e3", "user").Return("", errors.New("invalid role"))
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid json",
			requestBody:  "json",
			mock:         func(muu *MockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &MockUserUsecase{}
			tt.mock(mockUsecase)

			handler := delivery.NewRegisterHandler(mockUsecase)

			router := gin.Default()
			router.POST("/register", handler.Register)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "id")
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name         string
		requestBody  any
		mock         func(*MockUserUsecase)
		expectedCode int
	}{
		{
			name: "successful  login",
			requestBody: map[string]string{
				"email":    "kanzartem11@mail.ru",
				"password": "q1w2e3",
			},
			mock: func(muu *MockUserUsecase) {
				muu.On("Login", "kanzartem11@mail.ru", "q1w2e3").Return("token", nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "incorrect login",
			requestBody: map[string]string{
				"email":    "kanzartem1@mail.ru",
				"password": "q1w2e3",
			},
			mock: func(muu *MockUserUsecase) {
				muu.On("Login", "kanzartem1@mail.ru", "q1w2e3").Return("", errors.New("incorrect login"))
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "incorrect password",
			requestBody: map[string]string{
				"email":    "kanzartem11@mail.ru",
				"password": "q1w2",
			},
			mock: func(muu *MockUserUsecase) {
				muu.On("Login", "kanzartem11@mail.ru", "q1w2").Return("", errors.New("incorrect password"))
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "missing body",
			requestBody:  nil,
			mock:         func(muu *MockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid_json",
			requestBody:  "json",
			mock:         func(muu *MockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &MockUserUsecase{}
			tt.mock(mockUsecase)

			handler := delivery.NewLoginHandler(mockUsecase)

			router := gin.Default()
			router.POST("/login", handler.Login)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "token")
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}
