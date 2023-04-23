package auth

import "context"

type Storage interface {
	FindUserByEmail(ctx context.Context, email string) (*Profile, error)
	UpsertUser(ctx context.Context, user Profile) error
	DeleteUser(ctx context.Context, email string) error
}

type Profile struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
