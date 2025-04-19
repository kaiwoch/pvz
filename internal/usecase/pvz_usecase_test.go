package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"pvz/internal/storage/migrations/entity"
	"pvz/internal/usecase"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPVZStorage struct {
	mock.Mock
}

func (m *MockPVZStorage) CreatePVZ(id, user_id uuid.UUID, city string, date time.Time) (*entity.PVZ, error) {
	args := m.Called(id, user_id, city, date)
	return args.Get(0).(*entity.PVZ), args.Error(1)
}

func (m *MockPVZStorage) GetPVZById(id uuid.UUID) (*entity.PVZ, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.PVZ), args.Error(1)
}

func (m *MockPVZStorage) GetPVZsWithFilter(ctx context.Context, filter entity.Filter) ([]entity.ListPVZ, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]entity.ListPVZ), args.Error(1)
}

func (m *MockPVZStorage) CountPVZsWithFilter(ctx context.Context, filter entity.Filter) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

func TestPVZUsecase_CreatePVZ(t *testing.T) {
	pvz_id := uuid.Must(uuid.NewV4())
	userID := uuid.Must(uuid.NewV4())
	city := "Москва"
	date := time.Now()

	tests := []struct {
		name          string
		pvz_id        uuid.UUID
		user_id       uuid.UUID
		city          string
		date          time.Time
		expected      *entity.PVZ
		getPVZresult  *entity.PVZ
		expectedError error
		getPVZError   error
	}{
		{
			name:    "success",
			pvz_id:  pvz_id,
			user_id: userID,
			city:    city,
			expected: &entity.PVZ{
				ID:               pvz_id,
				RegistrationDate: date,
				City:             city,
				UserID:           userID,
			},
			expectedError: nil,
			getPVZError:   sql.ErrNoRows,
			getPVZresult:  &entity.PVZ{},
		},
		{
			name:          "pvz exists",
			pvz_id:        pvz_id,
			user_id:       userID,
			city:          city,
			expectedError: errors.New("pvz exists"),
			getPVZresult: &entity.PVZ{
				ID: pvz_id,
			},
			getPVZError: nil,
		},
		{
			name:          "db error",
			pvz_id:        pvz_id,
			user_id:       userID,
			city:          city,
			expectedError: errors.New("failed to check PVZ existence: db error"),
			getPVZError:   errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PVZStorage := new(MockPVZStorage)
			usecase := usecase.NewPVZUsecase(PVZStorage)

			PVZStorage.On("GetPVZById", tt.pvz_id).Return(tt.getPVZresult, tt.getPVZError)

			if tt.getPVZError == sql.ErrNoRows && tt.expectedError == nil {
				PVZStorage.On("CreatePVZ", tt.pvz_id, tt.user_id, tt.city, tt.date).Return(tt.expected, tt.expectedError)
			}

			pvz, err := usecase.CreatePVZ(tt.pvz_id, tt.user_id, tt.city, tt.date)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, pvz)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, pvz)
			}

			PVZStorage.AssertExpectations(t)
		})
	}
}

func TestPVZUsecase_GetPVZsWithFilter(t *testing.T) {
	pvz_id := uuid.Must(uuid.NewV4())
	reception_id := uuid.Must(uuid.NewV4())
	product_id := uuid.Must(uuid.NewV4())
	date := time.Now()

	filter := entity.Filter{
		StartDate: &date,
		EndDate:   &date,
		Page:      1,
		Limit:     10,
	}

	tests := []struct {
		name           string
		filter         entity.Filter
		getPVZresult   []entity.ListPVZ
		getPVZError    error
		countPVZresult int
		countPVZerror  error
		expectedError  error
		expected       *usecase.PVZListResponse
	}{
		{
			name:        "success",
			filter:      filter,
			getPVZError: nil,
			getPVZresult: []entity.ListPVZ{
				{
					Pvz: entity.PVZ{
						ID:               pvz_id,
						RegistrationDate: date,
						City:             "Москва",
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
			expectedError:  nil,
			countPVZresult: 1,
			countPVZerror:  nil,
			expected: &usecase.PVZListResponse{
				PVZs: []entity.ListPVZ{
					{
						Pvz: entity.PVZ{
							ID:               pvz_id,
							RegistrationDate: date,
							City:             "Москва",
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
				Total: 1,
				Page:  1,
				Limit: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PVZStorage := new(MockPVZStorage)
			usecase := usecase.NewPVZUsecase(PVZStorage)

			PVZStorage.On("GetPVZsWithFilter", context.Background(), tt.filter).Return(tt.getPVZresult, tt.getPVZError)

			if tt.getPVZError == nil && tt.expectedError == nil {
				PVZStorage.On("CountPVZsWithFilter", context.Background(), tt.filter).Return(tt.countPVZresult, tt.countPVZerror)
			}

			pvz_list, err := usecase.GetPVZsWithFilter(context.Background(), filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, pvz_list)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, pvz_list)
			}

			PVZStorage.AssertExpectations(t)
		})
	}
}
