package google

import (
	"time"

	"github.com/peoxia/auth-api/client"
	"github.com/peoxia/auth-api/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Client holds the HTTP client and Google APIs information.
type Client struct {
	OauthAPI    string
	PeopleAPI   string
	OauthConfig *oauth2.Config
	HTTPClient  client.HTTPClient
}

// Init sets up a new Google Oauth client.
func (c *Client) Init(config *config.Config) error {
	timeout := 1 * time.Second
	c.OauthAPI = "https://www.googleapis.com/oauth2/v2"
	c.PeopleAPI = "https://people.googleapis.com/v1"
	c.OauthConfig = &oauth2.Config{
		RedirectURL:  config.GoogleOauthCallbackURL,
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/user.phonenumbers.read",
		},
		Endpoint: google.Endpoint,
	}
	c.HTTPClient = client.NewHTTPClient(client.Parameters{Timeout: &timeout})
	return nil
}
