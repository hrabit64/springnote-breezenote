package auth

import (
	"context"
	"github.com/hrabit64/springnote-breezenote/config"
)

type FirebaseAuthClient struct {
}

func NewFirebaseAuthClient() *FirebaseAuthClient {
	return &FirebaseAuthClient{}
}

// VerifyIDToken firebase idToken을 검증하고, 해당 사용자의 uid를 반환합니다.
func (f *FirebaseAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*TokenInfo, error) {
	decodedToken, err := config.FirebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return &TokenInfo{}, err
	}

	return &TokenInfo{
		UID: decodedToken.UID,
	}, nil
}
