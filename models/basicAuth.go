package models

import (
	"context"
	"encoding/base64"
	"golang.org/x/oauth2"
)

type BasicAuth struct {
	Username string
	Password string
}

func (b BasicAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	auth := b.Username + ":" + b.Password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}

func (b BasicAuth) RequireTransportSecurity() bool {
	return true
}

func FetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}
