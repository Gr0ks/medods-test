package usecase

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"medods-test/pkg/models"
	"medods-test/pkg/auth"
	"time"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AccessPairCreator struct {
	repo auth.Repository

	hashSalt string
	signingKey []byte
	expireDuration time.Duration
}

func NewAccessPairCreator(repo auth.Repository, hashSalt string, signingKey []byte, expireDuration time.Duration) *AccessPairCreator {
	return &AccessPairCreator{
		repo: repo,
		hashSalt: hashSalt,
		signingKey: signingKey,
		expireDuration: expireDuration,
	}
}

func (a *AccessPairCreator) GetNewPair(session *models.Session) (*models.AccessPair, error) {
	startedAt := time.Now()
	session.StartedAt = startedAt.String()
	
	dataForRefreshKey := session.UserId + session.UserIP +session.StartedAt
	refreshKey :=sha1.New()
	refreshKey.Write([]byte(dataForRefreshKey))
	refreshKey.Write([]byte(a.hashSalt))
	refreshToken := fmt.Sprintf("%x", refreshKey.Sum(nil))
	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.MinCost)
	if err != nil {
		return &models.AccessPair{}, err
	}
	session.RefreshToken = string(refreshKeyHash)

	err := a.repo.Insert(session)
	if err != nil {
		return &models.AccessPair{}, err
	}

	expiredAt := startedAt.Add(a.expireDuration)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, &auth.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt.Unix(),
			IssuedAt: startedAt.Unix(),
		},
		UserId: session.UserId,
		UserIP: session.UserIP,
		StartedAt: session.StartedAt,
	})

	accessKey, err := accessToken.SignedString(a.signingKey)
	if err != nil {
		return &models.AccessPair{}, err
	}

	return &models.AccessPair{
		AccessKey: accessKey,
		RefreshKey: refreshToken,
	}, nil
}
	
func (a *AccessPairCreator) RefreshPair(accessPair *models.AccessPair, newSession *models.Session) (*models.AccessPair, error) {
	token, err := jwt.ParseWithClaims(accessPair.AccessKey, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return a.signingKey, nil
	})

	if err != nil {
		return &models.AccessPair{}, err
	}

	claims, ok := token.Claims.(*auth.Claims)
	if ok && token.Valid {
		session, err := a.repo.Get(claims.UserId)
		if err != nil {
			return &models.AccessPair{}, err
		}

		dataForRefreshKey := claims.UserId + claims.UserIP +claims.StartedAt
		refreshKey :=sha1.New()
		refreshKey.Write([]byte(dataForRefreshKey))
		refreshKey.Write([]byte(a.hashSalt))
		refreshToken := fmt.Sprintf("%x", refreshKey.Sum(nil))
		err := bcrypt.CompareHashAndPassword([]byte(session.RefreshToken), refreshToken)
		if err != nil {
			return &models.AccessPair{}, err
		}

		if refreshToken == accessPair.RefreshKey && refreshTokenHash == session.RefreshToken {
			return a.GetNewPair(newSession)
		}
		
		return &models.AccessPair{}, errors.New("invalid access token")
	}

	return &models.AccessPair{}, errors.New("invalid access token")
}