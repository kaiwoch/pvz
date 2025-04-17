package delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pvz/internal/delivery"
	"pvz/internal/storage/migrations/entity"
	"pvz/internal/usecase"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPVZUsecase struct {
	mock.Mock
}

func (p *MockPVZUsecase) CreatePVZ(id, user_id uuid.UUID, city string, date time.Time) (*entity.PVZ, error) {
	args := p.Called(id, user_id, city, date)
	return args.Get(0).(*entity.PVZ), args.Error(1)
}

func (p *MockPVZUsecase) GetPVZsWithFilter(ctx context.Context, filter entity.Filter) (*usecase.PVZListResponse, error) {
	args := p.Called(ctx, filter)
	return args.Get(0).(*usecase.PVZListResponse), args.Error(1)
}

func TestPostPVZHandler(t *testing.T) {
	pvz_id := uuid.Must(uuid.NewV4())
	user_id := uuid.Must(uuid.NewV4())
	date := time.Now()

	tests := []struct {
		name         string
		userID       string
		role         string
		requestBody  any
		mock         func(*MockPVZUsecase)
		expectedCode int
		expectedBody any
	}{
		{
			name:   "succesful created",
			role:   "moderator",
			userID: user_id.String(),
			requestBody: map[string]any{
				"id":               pvz_id,
				"registrationDate": date.Format(time.RFC3339),
				"city":             "Москва",
			},
			mock: func(mru *MockPVZUsecase) {
				mru.On("CreatePVZ", pvz_id, user_id, "Москва", mock.AnythingOfType("time.Time")).Return(&entity.PVZ{
					ID:               pvz_id,
					City:             "Москва",
					UserID:           user_id,
					RegistrationDate: date,
				}, nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:   "wrong role",
			role:   "employee",
			userID: user_id.String(),
			requestBody: map[string]any{
				"id":               pvz_id,
				"registrationDate": date.Format(time.RFC3339),
				"city":             "Москва",
			},
			mock:         func(mru *MockPVZUsecase) {},
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "invalid json",
			role:         "moderator",
			userID:       user_id.String(),
			requestBody:  "json",
			mock:         func(mru *MockPVZUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "empty body",
			role:         "moderator",
			userID:       user_id.String(),
			requestBody:  nil,
			mock:         func(mru *MockPVZUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "wrong city",
			role:   "moderator",
			userID: user_id.String(),
			requestBody: map[string]any{
				"id":               pvz_id,
				"registrationDate": date.Format(time.RFC3339),
				"city":             "Ступино",
			},
			mock: func(mru *MockPVZUsecase) {
				mru.On("CreatePVZ", pvz_id, user_id, "Ступино", mock.AnythingOfType("time.Time")).Return(&entity.PVZ{}, errors.New("bad city"))
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "bad user id",
			role:   "moderator",
			userID: "123",
			requestBody: map[string]any{
				"id":               pvz_id,
				"registrationDate": date.Format(time.RFC3339),
				"city":             "Москва",
			},
			mock:         func(mru *MockPVZUsecase) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &MockPVZUsecase{}
			tt.mock(mockUsecase)

			handler := delivery.NewPVZHandler(mockUsecase)

			router := gin.Default()
			router.POST("/pvz", func(ctx *gin.Context) {
				ctx.Set("userID", tt.userID)
				ctx.Set("role", tt.role)
				handler.PostPVZ(ctx)
			})

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestGetPVZHandler(t *testing.T) {
	var startDate, endDate time.Time
	startDate = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	pvz_id := uuid.Must(uuid.NewV4())
	date := time.Now()
	user_id := uuid.Must(uuid.NewV4())
	reception_id := uuid.Must(uuid.NewV4())
	product_id := uuid.Must(uuid.NewV4())

	tests := []struct {
		name         string
		role         string
		queryParams  string
		requestBody  any
		mock         func(*MockPVZUsecase)
		expectedCode int
		expectedBody any
	}{
		{
			name: "succesful get",
			role: "moderator",
			queryParams: fmt.Sprintf("startDate=%s&endDate=%s&page=1&limit=10",
				startDate.Format(time.RFC3339),
				endDate.Format(time.RFC3339)),
			requestBody: nil,
			mock: func(mru *MockPVZUsecase) {
				mru.On("GetPVZsWithFilter", context.Background(), entity.Filter{
					StartDate: &startDate,
					EndDate:   &endDate,
					Page:      1,
					Limit:     10,
				}).Return(&usecase.PVZListResponse{
					PVZs: []entity.ListPVZ{
						{
							Pvz: entity.PVZ{
								ID:               pvz_id,
								City:             "Москва",
								UserID:           user_id,
								RegistrationDate: date,
							},
							Receptions: []entity.Receptions{
								{
									ID:       reception_id,
									DateTime: date,
									PVZID:    pvz_id,
									Status:   "close",
									Products: []entity.Products{
										{
											ID:          product_id,
											DateTime:    date,
											Type:        "одежда",
											ReceptionId: reception_id,
										},
									},
								},
							},
						},
					},
					Page:  1,
					Limit: 10,
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]any{
				"pvzs": []any{
					map[string]any{
						"pvz": map[string]any{
							"UserID":            user_id.String(),
							"city_name":         "Москва",
							"pvz_id":            pvz_id.String(),
							"registration_date": date.Format(time.RFC3339Nano),
						},
						"receptions": []any{
							map[string]any{
								"dateTime": date.Format(time.RFC3339Nano),
								"id":       reception_id.String(),
								"products": []any{
									map[string]any{
										"dateTime":    date.Format(time.RFC3339Nano),
										"id":          product_id.String(),
										"receptionId": reception_id.String(),
										"type":        "одежда",
									},
								},
								"pvzId":  pvz_id.String(),
								"status": "close",
							},
						},
					},
				},
				"page":  float64(1),
				"limit": float64(10),
				"total": float64(0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &MockPVZUsecase{}
			tt.mock(mockUsecase)

			handler := delivery.NewPVZHandler(mockUsecase)

			router := gin.Default()
			router.GET("/pvz", func(ctx *gin.Context) {
				ctx.Set("role", tt.role)
				handler.GetPVZs(ctx)
			})

			req, _ := http.NewRequest(http.MethodGet, "/pvz?"+tt.queryParams, nil)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedBody != nil {
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
			mockUsecase.AssertExpectations(t)
		})
	}
}
