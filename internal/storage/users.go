package storage

import (
	"database/sql"
	"pvz/internal/storage/migrations/entity"

	"github.com/gofrs/uuid/v5"
)

type UsersPostgresStorage interface {
	GetUserByEmail(email string) (*entity.User, bool, error)
	CreateUser(email, password, role string) (*entity.User, error)
}

type UsersPostgresStorageImpl struct {
	db *sql.DB
}

func NewUsersStorage(db *sql.DB) *UsersPostgresStorageImpl {
	return &UsersPostgresStorageImpl{db: db}
}

func (u *UsersPostgresStorageImpl) GetUserByEmail(email string) (*entity.User, bool, error) {
	var user entity.User
	query := "SELECT * FROM users WHERE email = $1"
	err := u.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil && err != sql.ErrNoRows {
		return nil, false, err
	} else if err == sql.ErrNoRows {
		return nil, false, sql.ErrNoRows
	}

	return &user, true, nil
}

func (u *UsersPostgresStorageImpl) CreateUser(email, password, role string) (*entity.User, error) {
	id := uuid.Must(uuid.NewV4())
	query := "INSERT INTO users (user_id, email, password_hash, role_name) VALUES ($1, $2, $3, $4)"
	_, err := u.db.Exec(query, id, email, password, role)
	if err != nil {
		return nil, err
	}
	return &entity.User{ID: id, Email: email, Password: string(password), Role: role}, nil
}
