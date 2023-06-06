package postgresql

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) CreateUser(ctx context.Context, number, password, name, lastname string) (err error) {
	_, err = r.db.ExecContext(ctx, `INSERT INTO users (number, password, name, lastname) VALUES ($1, $2, $3, $4)`, number, password, name, lastname)
	return
}

func (r UserRepository) GetUser(ctx context.Context, number, password string) (Id_users int64, err error) {

	err = r.db.Get(&Id_users, "SELECT id_users FROM users WHERE number=$1 AND password=$2", number, password)
	return
}
