package mock

import (
	"context"
	"errors"
	"github.com/hrabit64/springnote-breezenote/auth"
	"github.com/hrabit64/springnote-breezenote/config"
)

var (
	// AdminToken 관리자용 토큰
	AdminToken = "imadmin"

	// UserToken 일반 사용자용 토큰
	UserToken = "imuser"
)

type MockAuthClient struct{}

func NewMockAuthClient() *MockAuthClient {
	return &MockAuthClient{}
}

func (m *MockAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.TokenInfo, error) {
	switch {
	case idToken == AdminToken:
		return &auth.TokenInfo{
			UID: config.RootConfig.AllowImageUploadUid,
		}, nil

	case idToken == UserToken:
		return &auth.TokenInfo{
			UID: "user",
		}, nil

	default:
		return &auth.TokenInfo{}, errors.New("invalid token")
	}

}
