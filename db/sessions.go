package db

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

type CreateSessionRequest struct {
	Id           uuid.UUID
	UserId       int
	RefreshToken string
	ClientIp     string
	UserAgent    string
	ExpiresAt    time.Time
}

func (d *DBStruct) CreateSession(ctx context.Context, body CreateSessionRequest) (Session, error) {
	stmt, err := d.DB.PrepareContext(ctx, "INSERT INTO sessions (id, user_id, refresh_token, client_ip, user_agent, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		return Session{}, err
	}

	createdAt := time.Now()
	result, err := stmt.ExecContext(ctx, body.Id, body.UserId, body.RefreshToken, body.ClientIp, body.UserAgent, body.ExpiresAt, createdAt)

	if err != nil {
		return Session{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Session{}, err
	}
	log.Println("ID Here", id)

	return Session{Id: body.Id, UserId: body.UserId, RefreshToken: body.RefreshToken, ClientIp: body.ClientIp, UserAgent: body.UserAgent, ExpiresAt: body.ExpiresAt, CreatedAt: createdAt}, nil
}

func (d *DBStruct) GetSession(ctx context.Context, id uuid.UUID) (Session, error) {
	var session Session
	row := d.DB.QueryRowContext(ctx, "SELECT id, user_id, refresh_token, client_ip, user_agent, expires_at, created_at from sessions WHERE id = ?;", id)
	err := row.Scan(&session.Id, &session.UserId, &session.RefreshToken, &session.ClientIp, &session.UserAgent, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return Session{}, err
	}
	return session, nil
}
