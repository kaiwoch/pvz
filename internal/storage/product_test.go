package storage_test

import (
	"pvz/internal/storage"
	"pvz/internal/storage/migrations/entity"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
)

func TestProductPostgresStorage_CreateProduct(t *testing.T) {
	reception_id := uuid.Must(uuid.NewV4())
	product_type := "одежда"
	date := time.Now()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := storage.NewProductPostgresStorage(db)

	tests := []struct {
		name         string
		reception_id uuid.UUID
		product_type string
		mock         func()
		expected     *entity.Products
		expectedErr  error
	}{
		{
			name:         "success",
			reception_id: reception_id,
			product_type: product_type,
			mock: func() {
				mock.ExpectExec("INSERT INTO product").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), product_type, reception_id).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expected: &entity.Products{
				DateTime:    date,
				Type:        product_type,
				ReceptionId: reception_id,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			product, err := storage.CreateProduct(tt.reception_id, tt.product_type)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Nil(t, product)
			} else {
				assert.NotEqual(t, uuid.Nil, product.ID)
				assert.NoError(t, err)
				assert.NotNil(t, product)
				assert.Equal(t, tt.expected.ReceptionId, product.ReceptionId)
				assert.WithinDuration(t, time.Now(), product.DateTime, time.Second)
				assert.Equal(t, tt.expected.Type, product.Type)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProductPostgresStorage_GetLastProductID(t *testing.T) {
	reception_id := uuid.Must(uuid.NewV4())
	product_id := uuid.Must(uuid.NewV4())

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := storage.NewProductPostgresStorage(db)

	tests := []struct {
		name         string
		reception_id uuid.UUID
		product_id   uuid.UUID
		mock         func()
		expected     uuid.UUID
		expectedErr  error
	}{
		{
			name:         "success",
			reception_id: reception_id,
			mock: func() {
				rows := mock.NewRows([]string{"product_id"}).AddRow(product_id)
				mock.ExpectQuery("SELECT product_id FROM product WHERE reception_id = \\$1 ORDER BY date_time DESC LIMIT 1").
					WithArgs(reception_id).WillReturnRows(rows)
			},
			expected:    product_id,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			product_id, err := storage.GetLastProductID(tt.reception_id)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Nil(t, product_id)
			} else {
				assert.NotEqual(t, uuid.Nil, product_id)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, product_id)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
