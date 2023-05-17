package auth

import "time"

type Maker interface {
	CreateToken(userId int, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
