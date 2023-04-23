package auth

import "golang.org/x/oauth2"

type LoginProvider interface {
	GetAuthCodeURL() (*AuthCodeURL, error)
	GetAccessToken(accessCode string) (*oauth2.Token, error)
	GetUserInfo(accessToken oauth2.Token) (*Profile, error)
}

type AuthCodeURL struct {
	URL   string
	State string
}
