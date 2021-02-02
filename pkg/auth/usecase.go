package auth

import (
	"medods-test/pkg/models"
	"context"
)

type UseCase interface {
	GetNewPair(ctx context.Context, session *models.Session) (*models.AccessPair, error)
	RefreshPair(ctx context.Context, accessPair *models.AccessPair, newSession *models.Session) (*models.AccessPair, error)
}