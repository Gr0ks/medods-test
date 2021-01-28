package mongo

import (
	"context"
	"medods-test/pkg/models"
	"medods-test/pkg/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"errors"
)

type SessionRepository struct {
	db *mongo.Collection
}

func NewSessionRepository(db *mongo.Database, collection string) *SessionRepository {
	return &SessionRepository{
		db: db.Collection(collection),
	}
}

func (r *SessionRepository) InsertOrUpdate(session *models.Session) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"userId", session.UserId}}
	update := bson.D{{"$set", bson.D{
		{"userIP", session.UserIP},
		{"startedAt", session.StartedAt},
		{"refreshToken", session.RefreshToken},
	}}}
	err := r.db.FindOneAndUpdate(filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func (r *SessionRepository) Get(userId string) (*models.Session, error) {
	session := new(models.Session)

	if err := r.db.FindOne(ctx, bson.M{"userId": userId}).Decode(session); err != nil {
		log.Errorf("error occured while getting user from db: %s", err.Error())
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user does not exist")
		}

		return nil, err
	}

	return session, nil
}