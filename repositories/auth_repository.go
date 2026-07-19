package repositories

import (
	"database/sql"
	"errors"
	"management-stock/models"
)

type AuthRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(query, user.Email, user.Password, user.Role).Scan(&user.ID, &user.CreatedAt)
}

func (r *authRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password, role, created_at FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
