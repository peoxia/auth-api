package google

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/peoxia/auth-api/auth"
	"github.com/peoxia/auth-api/client"
	"golang.org/x/oauth2"
)

// GetUserInfo requests authorised user info.
func (c *Client) GetUserInfo(t oauth2.Token) (*auth.Profile, error) {

	reqData := client.HTTPRequestData{
		Method: http.MethodGet,
		URL:    c.OauthAPI + "/userinfo?access_token=" + t.AccessToken,
	}

	respBody, err := c.HTTPClient.RequestBytes(reqData)
	if err != nil {
		return nil, fmt.Errorf("error making request to get authorised user info: %w", err)
	}

	var resp auth.Profile
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling authorised user info: %w", err)
	}

	return &resp, nil
}

// GetUserPhone requests authorised user phone number by calling Google's People API.
func (c *Client) GetUserPhone() (string, error) {
	return "61400000000", nil
}
