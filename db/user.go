package db

import "context"

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

func (db *DBStruct) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	row := db.DB.QueryRowContext(ctx, "SELECT id, first_name, last_name, password, email, created_at, updated_at from users WHERE email = ?", email)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}
