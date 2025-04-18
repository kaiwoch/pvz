package storage_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid/v5"

	"github.com/stretchr/testify/assert"

	"pvz/internal/storage"
	"pvz/internal/storage/migrations/entity"
)

func TestUsersStorage_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := storage.NewUsersStorage(db)

	tests := []struct {
		name        string
		email       string
		mock        func()
		expected    *entity.User
		expectedErr error
	}{
		{
			name:  "success",
			email: "test@example.com",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "email", "password_hash", "role_name"}).
					AddRow("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "test@example.com", "hash123", "moderator")
				mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\$1").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expected: &entity.User{
				ID:       uuid.FromStringOrNil("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				Email:    "test@example.com",
				Password: "hash123",
				Role:     "moderator",
			},
			expectedErr: nil,
		},
		{
			name:  "not found",
			email: "test@example.com",
			mock: func() {
				mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\$1").
					WithArgs("test@example.com").
					WillReturnError(fmt.Errorf("user not exists"))
			},
			expectedErr: fmt.Errorf("user not exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			user, _, err := storage.GetUserByEmail(tt.email)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, user)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUsersStorage_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := storage.NewUsersStorage(db)

	tests := []struct {
		name        string
		email       string
		password    string
		role        string
		mock        func()
		expectedErr error
	}{
		{
			name:     "success",
			email:    "test@example.com",
			password: "hash123",
			role:     "moderator",
			mock: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(sqlmock.AnyArg(), "test@example.com", "hash123", "moderator").WillReturnResult(sqlmock.NewResult(1, 1))

			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			user, err := storage.CreateUser(tt.email, tt.password, tt.role)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.password, user.Password)
				assert.Equal(t, tt.role, user.Role)
				assert.NotEqual(t, uuid.Nil, user.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
