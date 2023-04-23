package google

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/peoxia/auth-api/auth"
	"golang.org/x/oauth2"
)

// GetAuthCodeURL creates a state and returns auth code URL.
func (c *Client) GetAuthCodeURL() (*auth.AuthCodeURL, error) {
	state, err := randomString(16)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate a state")
	}
	return &auth.AuthCodeURL{
		URL:   c.OauthConfig.AuthCodeURL(state),
		State: state,
	}, nil
}

// GetAccessToken exchanges access code for access token.
func (c *Client) GetAccessToken(accessCode string) (*oauth2.Token, error) {
	return c.OauthConfig.Exchange(context.Background(), accessCode)
}

func randomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
