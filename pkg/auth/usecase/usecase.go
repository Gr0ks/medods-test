package usecase

import (
	"crypto/sha1"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"medods-test/pkg/models"
	"medods-test/pkg/auth"
	"time"
	"errors"
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
	
	refreshKey := sha1.New()
	refreshKey.Write([]byte(session.UserId))
	refreshKey.Write([]byte(session.UserIP))
	refreshKey.Write([]byte(session.StartedAt))
	refreshKey.Write([]byte(a.hashSalt))
	refreshToken := fmt.Sprintf("%x", refreshKey.Sum(nil))

	savedRefreshToken := sha1.New()
	savedRefreshToken.Write([]byte(refreshToken))
	savedRefreshToken.Write([]byte(a.hashSalt))

	session.RefreshToken = fmt.Sprintf("%x", savedRefreshToken.Sum(nil))

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
		RefreshKey: session.RefreshToken,
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

		refreshKey :=sha1.New()
		refreshKey.Write([]byte(claims.UserId))
		refreshKey.Write([]byte(claims.UserIP))
		refreshKey.Write([]byte(claims.StartedAt))
		refreshKey.Write([]byte(a.hashSalt))
		refreshToken := fmt.Sprintf("%x", refreshKey.Sum(nil))

		refreshTokenKey := sha1.New()
		refreshTokenKey.Write([]byte(accessPair.RefreshKey))
		refreshTokenKey.Write([]byte(a.hashSalt))

		refreshTokenHash := fmt.Sprintf("%x", refreshTokenKey.Sum(nil))

		if refreshToken == accessPair.RefreshKey && refreshTokenHash == session.RefreshToken {
			return a.GetNewPair(newSession)
		}
		
		return &models.AccessPair{}, errors.New("invalid access token")
	}

	return &models.AccessPair{}, errors.New("invalid access token")
}