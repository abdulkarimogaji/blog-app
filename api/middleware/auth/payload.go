package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID
	UserID    int
	IssuedAt  time.Time
	ExpiredAt time.Time
}

var (
	ErrTokenExpired = errors.New("TOKEN_EXPIRED")
	ErrInValidToken = errors.New("INVALID_TOKEN")
)

func NewPayload(userId int, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:        tokenId,
		UserID:    userId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrTokenExpired
	}
	return nil
}
