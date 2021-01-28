package auth

import (
	"medods-test/pkg/models"
)

type Repository interface {
	Insert(session *models.Session) error
	Get(userId string) (*models.Session, error)
}