package auth

import (
	"context"
)

type AuthenticateClient interface {
	VerifyIDToken(ctx context.Context, idToken string) (*TokenInfo, error)
}

type TokenInfo struct {
	UID string
}
