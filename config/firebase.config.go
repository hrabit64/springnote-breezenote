package config

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var FirebaseAuth *auth.Client

// SetFirebaseAuth firebase 인증을 설정합니다. 루트 디렉토리에 firebase.json 파일이 있어야 합니다.
func SetFirebaseAuth() error {
	opt := option.WithCredentialsFile(RootConfig.FireBaseConfig)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	FirebaseAuth, err = app.Auth(context.Background())
	if err != nil {
		return err
	}

	return nil
}
