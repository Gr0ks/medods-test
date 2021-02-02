package mongo

import (
	"context"
	"medods-test/pkg/models"
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

func (r *SessionRepository) InsertOrUpdate(ctx context.Context, session *models.Session) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"userId", session.UserId}}
	update := bson.D{{"$set", bson.D{
		{"userIP", session.UserIP},
		{"startedAt", session.StartedAt},
		{"refreshToken", session.RefreshToken},
	}}}
	r.db.FindOneAndUpdate(ctx, filter, update, opts)
	return nil
}

func (r *SessionRepository) Get(ctx context.Context, userId string) (*models.Session, error) {
	session := new(models.Session)
	if err := r.db.FindOne(ctx, bson.M{"userId": userId}).Decode(session); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user does not exist")
		}
		return nil, err
	}
	return session, nil
}