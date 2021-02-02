package auth

import (
	"medods-test/pkg/models"
	"context"
)

type Repository interface {
	InsertOrUpdate(ctx context.Context, session *models.Session) error
	Get(ctx context.Context, userId string) (*models.Session, error)
}