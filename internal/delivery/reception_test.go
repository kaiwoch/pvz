package delivery_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pvz/internal/delivery"
	"pvz/internal/storage/migrations/entity"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReceptionUsecase struct {
	mock.Mock
}

func (m *MockReceptionUsecase) CreateReception(id uuid.UUID) (*entity.Receptions, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Receptions), args.Error(1)
}

func (m *MockReceptionUsecase) UpdateReceptionStatus(pvz_id uuid.UUID) (*entity.Receptions, error) {
	args := m.Called(pvz_id)
	return args.Get(0).(*entity.Receptions), args.Error(1)
}

func TestReceptionHandler(t *testing.T) {
	receptionID := uuid.Must(uuid.NewV4())
	pvzID := uuid.Must(uuid.NewV4())
	date := time.Now()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		role         string
		requestBody  any
		mock         func(*MockReceptionUsecase)
		expectedCode int
		expectedBody any
	}{
		{
			name: "successfully created reception",
			role: "employee",
			requestBody: map[string]any{
				"pvzId": pvzID.String(),
			},
			mock: func(m *MockReceptionUsecase) {
				m.On("CreateReception", pvzID).Return(&entity.Receptions{
					ID:       receptionID,
					DateTime: date,
					PVZID:    pvzID,
					Status:   "in_progress",
				}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: gin.H{
				"id":       receptionID.String(),
				"dateTime": date.Format(time.RFC3339),
				"pvzId":    pvzID.String(),
				"status":   "in_progress",
			},
		},
		{
			name: "wrong role",
			role: "moderator",
			requestBody: map[string]any{
				"pvzId": pvzID.String(),
			},
			mock:         func(m *MockReceptionUsecase) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{},
		},
		{
			name:         "invalid json",
			role:         "employee",
			requestBody:  "json",
			mock:         func(m *MockReceptionUsecase) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{},
		},
		{
			name:         "empty body",
			role:         "employee",
			requestBody:  nil,
			mock:         func(m *MockReceptionUsecase) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{},
		},
		{
			name: "previous reception is not closed",
			role: "employee",
			requestBody: map[string]any{
				"pvzId": pvzID.String(),
			},
			mock: func(m *MockReceptionUsecase) {
				m.On("CreateReception", pvzID).Return(&entity.Receptions{}, errors.New("no available receptions"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &MockReceptionUsecase{}
			tt.mock(mockUsecase)

			handler := delivery.NewReceptionHandler(mockUsecase)

			router := gin.Default()
			router.POST("/receptions", func(ctx *gin.Context) {
				ctx.Set("role", tt.role)
				handler.Reception(ctx)
			})

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/receptions", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusCreated {
				var response map[string]any
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody.(gin.H)["id"], response["id"])
				assert.Equal(t, tt.expectedBody.(gin.H)["status"], response["status"])
				assert.Equal(t, tt.expectedBody.(gin.H)["pvzId"], response["pvzId"])

				respTime, err := time.Parse(time.RFC3339, response["dateTime"].(string))
				assert.NoError(t, err)
				assert.WithinDuration(t, date, respTime, time.Second)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUpdateReceptionHandler(t *testing.T) {
	receptionID := uuid.Must(uuid.NewV4())
	pvzID := uuid.Must(uuid.NewV4())
	date := time.Now()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		param        uuid.UUID
		name         string
		role         string
		requestBody  any
		mock         func(*MockReceptionUsecase)
		expectedCode int
		expectedBody any
	}{
		{
			name:        "successfully closed reception",
			param:       pvzID,
			role:        "employee",
			requestBody: nil,
			mock: func(m *MockReceptionUsecase) {
				m.On("UpdateReceptionStatus", pvzID).Return(&entity.Receptions{
					ID:       receptionID,
					DateTime: date,
					PVZID:    pvzID,
					Status:   "closed",
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{
				"id":       receptionID.String(),
				"dateTime": date.Format(time.RFC3339),
				"pvzId":    pvzID.String(),
				"status":   "closed",
			},
		},
		{
			name:         "wrong role",
			param:        pvzID,
			role:         "moderator",
			requestBody:  nil,
			mock:         func(m *MockReceptionUsecase) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &MockReceptionUsecase{}
			tt.mock(mockUsecase)

			handler := delivery.NewReceptionHandler(mockUsecase)

			router := gin.Default()
			router.POST("/pvz/:pvzId/close_last_reception", func(ctx *gin.Context) {
				ctx.Set("role", tt.role)
				handler.UpdateReceptionStatus(ctx)
			})

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/pvz/%s/close_last_reception", tt.param), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response map[string]any
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				fmt.Println(response)

				assert.Equal(t, tt.expectedBody.(gin.H)["id"], response["id"])
				assert.Equal(t, tt.expectedBody.(gin.H)["status"], response["status"])
				assert.Equal(t, tt.expectedBody.(gin.H)["pvzId"], response["pvzId"])

				respTime, err := time.Parse(time.RFC3339, response["dateTime"].(string))
				assert.NoError(t, err)
				assert.WithinDuration(t, date, respTime, time.Second)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}
