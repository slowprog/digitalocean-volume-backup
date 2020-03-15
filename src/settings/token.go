package settings

import (
	"golang.org/x/oauth2"
)

// TokenSource struct with oauth-2token
type TokenSource struct {
	AccessToken string
}

// Token implements interface of oauth2
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}

	return token, nil
}

// NewTokenSource returns a new oauth2-struct
func NewTokenSource(accessToken string) *TokenSource {
	return &TokenSource{
		AccessToken: accessToken,
	}
}
