package usecase_test

import (
	"database/sql"
	"errors"
	"pvz/internal/storage/migrations/entity"
	"pvz/internal/usecase"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReceptionStorage struct {
	mock.Mock
}

func (m *MockReceptionStorage) CreateReception(id uuid.UUID) (*entity.Receptions, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Receptions), args.Error(1)
}

func (m *MockReceptionStorage) GetLastReceptionStatus(id uuid.UUID) (uuid.UUID, string, error) {
	args := m.Called(id)
	return args.Get(0).(uuid.UUID), args.String(1), args.Error(2)
}

func (m *MockReceptionStorage) UpdateReceptionStatus(reception_id uuid.UUID) error {
	args := m.Called(reception_id)
	return args.Error(0)
}

func (m *MockReceptionStorage) GetReceptionById(reception_id uuid.UUID) (*entity.Receptions, error) {
	args := m.Called(reception_id)
	return args.Get(0).(*entity.Receptions), args.Error(1)
}

func TestReceptionUsecase_CreateReception(t *testing.T) {
	pvz_id := uuid.Must(uuid.NewV4())

	tests := []struct {
		name               string
		pvz_id             uuid.UUID
		getReceptionresult string
		getReceptionError  error
		expected           *entity.Receptions
		expectedError      error
	}{
		{
			name:               "success",
			pvz_id:             pvz_id,
			getReceptionError:  nil,
			getReceptionresult: "close",
			expected: &entity.Receptions{
				Status: "in_progress",
				PVZID:  pvz_id,
			},
			expectedError: nil,
		},
		{
			name:               "status check error",
			pvz_id:             pvz_id,
			getReceptionError:  sql.ErrNoRows,
			getReceptionresult: "",
			expected:           nil,
			expectedError:      errors.New("failed to check reception status: sql: no rows in result set"),
		},
		{
			name:               "status in_progress",
			pvz_id:             pvz_id,
			getReceptionError:  nil,
			getReceptionresult: "in_progress",
			expected:           nil,
			expectedError:      errors.New("close previous receipt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReceptionStorage := new(MockReceptionStorage)
			usecase := usecase.NewReceptionUsecase(ReceptionStorage)

			ReceptionStorage.On("GetLastReceptionStatus", tt.pvz_id).Return(uuid.UUID{}, tt.getReceptionresult, tt.getReceptionError)

			if tt.getReceptionError == nil && tt.expectedError == nil && tt.getReceptionresult == "close" {
				ReceptionStorage.On("CreateReception", tt.pvz_id).Return(tt.expected, tt.expectedError)
			}

			reception, err := usecase.CreateReception(tt.pvz_id)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, reception)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, reception)
			}
			ReceptionStorage.AssertExpectations(t)
		})
	}
}

func TestReceptionUsecase_UpdateReceptionStatus(t *testing.T) {
	pvz_id := uuid.Must(uuid.NewV4())

	tests := []struct {
		name               string
		pvz_id             uuid.UUID
		getReceptionResult struct {
			reception_id uuid.UUID
			status       string
		}
		getReceptionError    error
		updateReceptionError error
		expected             *entity.Receptions
		expectedError        error
	}{
		{
			name:              "success",
			pvz_id:            pvz_id,
			getReceptionError: nil,
			getReceptionResult: struct {
				reception_id uuid.UUID
				status       string
			}{
				reception_id: uuid.FromStringOrNil("3fa85f64-5717-4562-b3fc-2c963f66afa6"),
				status:       "in_progress",
			},
			expected: &entity.Receptions{
				ID:     uuid.FromStringOrNil("3fa85f64-5717-4562-b3fc-2c963f66afa6"),
				Status: "close",
				PVZID:  pvz_id,
			},
			expectedError:        nil,
			updateReceptionError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReceptionStorage := new(MockReceptionStorage)
			usecase := usecase.NewReceptionUsecase(ReceptionStorage)

			ReceptionStorage.On("GetLastReceptionStatus", tt.pvz_id).Return(tt.getReceptionResult.reception_id, tt.getReceptionResult.status, tt.getReceptionError)

			if tt.getReceptionError == nil && tt.expectedError == nil && tt.getReceptionResult.status == "in_progress" {
				ReceptionStorage.On("UpdateReceptionStatus", tt.getReceptionResult.reception_id).Return(tt.updateReceptionError)
				if tt.updateReceptionError == nil {
					ReceptionStorage.On("GetReceptionById", tt.getReceptionResult.reception_id).Return(tt.expected, tt.expectedError)
				}

			}

			reception, err := usecase.UpdateReceptionStatus(tt.pvz_id)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, reception)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, reception)
			}
			ReceptionStorage.AssertExpectations(t)
		})
	}
}
