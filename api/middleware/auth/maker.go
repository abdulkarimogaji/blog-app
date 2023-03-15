package auth

import "time"

type Maker interface {
	CreateToken(userId int, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
