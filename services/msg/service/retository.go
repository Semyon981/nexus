package service

import (
	"context"
	"time"

	"github.com/Semyon981/nexus/services/msg/models"
)

type Repository interface {
	SendMessage(ctx context.Context, Id_from int64, id_to int64, Msg string, Time time.Time) error
	GetMessages(ctx context.Context, Id_from int64, id_to int64, limit int64, offset int64) ([]models.Message, error)
}
