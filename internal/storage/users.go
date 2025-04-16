package storage

import (
	"database/sql"
	"pvz/internal/storage/migrations/entity"

	"github.com/gofrs/uuid/v5"
)

type UsersPostgresStorage struct {
	db *sql.DB
}

func NewUsersStorage(db *sql.DB) *UsersPostgresStorage {
	return &UsersPostgresStorage{db: db}
}

func (u *UsersPostgresStorage) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	query := "SELECT * FROM users WHERE email = $1"
	err := u.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UsersPostgresStorage) CreateUser(email, password, role string) (*entity.User, error) {
	id := uuid.Must(uuid.NewV4())
	query := "INSERT INTO users (user_id, email, password_hash, role_name) VALUES ($1, $2, $3, $4)"
	_, err := u.db.Exec(query, id, email, password, role)
	if err != nil {
		return nil, err
	}
	return &entity.User{ID: id, Email: email, Password: string(password), Role: role}, nil
}
