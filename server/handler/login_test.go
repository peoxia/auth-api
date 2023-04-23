package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/peoxia/auth-api/auth"
	"github.com/peoxia/auth-api/config"
	"golang.org/x/oauth2"
)

func TestLoginCallback(t *testing.T) {
	testCases := []struct {
		name             string
		cookie           *http.Cookie
		query            string
		provider         *mockLoginProvider
		storage          *mockStorage
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "successful login",
			cookie: &http.Cookie{
				Name:  "oauth-state",
				Value: "my-state",
			},
			query: "code=my-code&state=my-state",
			provider: &mockLoginProvider{
				accessToken: &oauth2.Token{
					AccessToken: "my-token",
					Expiry:      time.Now().Add(1 * time.Hour),
				},
				userInfo: &auth.Profile{
					ID:    "my-id",
					Name:  "my-name",
					Email: "my-email",
				},
			},
			storage:          &mockStorage{},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"id":"my-id","email":"my-email","name":"my-name"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/login/callback?"+tc.query, nil)
			req.AddCookie(tc.cookie)

			w := httptest.NewRecorder()
			
			config := config.Config{JWTSecret: "TEST"}
			LoginCallback(config, tc.provider, tc.storage)(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d but got %d", tc.expectedStatus, w.Code)
			}
			if w.Body.String() != tc.expectedResponse {
				t.Errorf("expected response %q but got %q", tc.expectedResponse, w.Body.String())
			}
		})
	}
}

type mockLoginProvider struct {
	accessToken *oauth2.Token
	userInfo    *auth.Profile
}

func (m *mockLoginProvider) GetAccessToken(code string) (*oauth2.Token, error) {
	return m.accessToken, nil
}

func (m *mockLoginProvider) GetUserInfo(accessToken oauth2.Token) (*auth.Profile, error) {
	return m.userInfo, nil
}
func (m *mockLoginProvider) GetAuthCodeURL() (*auth.AuthCodeURL, error) {
	return &auth.AuthCodeURL{}, nil
}

type mockStorage struct{}

func (m *mockStorage) UpsertUser(ctx context.Context, user auth.Profile) error {
	return nil
}
func (m *mockStorage) FindUserByEmail(ctx context.Context, email string) (*auth.Profile, error) {
	return &auth.Profile{}, nil
}
func (m *mockStorage) DeleteUser(ctx context.Context, email string) error {
	return nil
}
