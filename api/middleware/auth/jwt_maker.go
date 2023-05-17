package auth

import (
	"errors"
	"time"

	"github.com/abdulkarimogaji/blognado/config"
	"github.com/golang-jwt/jwt/v4"
)

type JwtMaker struct {
	secretKey string
}

func NewJwtMaker() (Maker, error) {
	return &JwtMaker{
		secretKey: config.AppConfig.JWT_SECRET,
	}, nil
}

func (m *JwtMaker) CreateToken(userId int, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userId, duration)
	if err != nil {
		return "", payload, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(m.secretKey))
	return tokenString, payload, err
}
func (m *JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInValidToken
		}
		return []byte(m.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		if v_err, ok := err.(*jwt.ValidationError); ok && errors.Is(v_err.Inner, ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInValidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, ErrInValidToken
	}
	return payload, nil
}
