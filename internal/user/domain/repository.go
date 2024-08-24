package domain

import "context"

type UserRepository interface {
	GetUser(ctx context.Context, username string) (*User, error)
}
