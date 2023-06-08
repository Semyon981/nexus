package postgresql

import (
	"context"
	"time"

	"github.com/Semyon981/nexus/services/msg/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) SendMessage(ctx context.Context, Id_from int64, Id_to int64, Msg string, Time time.Time) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO messages (id_from, id_to, msg, time) VALUES ($1, $2, $3, $4)`, Id_from, Id_to, Msg, Time)
	return err
}

func (r Repository) GetMessages(ctx context.Context, Id_from int64, Id_to int64, limit int64, offset int64) ([]models.Message, error) {
	res := []models.Message{}
	err := r.db.Select(&res, `SELECT * FROM messages WHERE (id_from = $3 AND id_to = $4) OR (id_from = $4 AND id_to = $3) ORDER BY time DESC LIMIT $1 OFFSET $2 `, limit, offset, Id_from, Id_to)
	return res, err
}
