package auth

import (
	"medods-test/pkg/models"
)

type UseCase interface {
	GetNewPair(session *models.Session) (*models.AccessPair, error)
	RefreshPair(accessPair *models.AccessPair, newSession *models.Session) (*models.AccessPair, error)
}