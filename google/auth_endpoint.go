package google

import (
	"net/url"
)

func (c *Client) BuildAuthURL() *url.URL {
	opt := c.option
	secret := c.secret.Installed

	authEndpointURL, err := url.Parse(opt.AuthEndpoint)
	if err != nil {
		panic(err)
	}

	q := authEndpointURL.Query()
	q.Add("client_id", secret.ClientID)
	q.Add("state", opt.State)
	q.Add("redirect_uri", opt.RedirectURI)
	q.Add("response_type", opt.ResponseType)
	q.Add("scope", opt.Scope)

	// PKCE用パラメータ
	q.Add("code_challenge", opt.CodeChallenge)
	q.Add("code_challenge_method", opt.CodeChallengeMethod)

	authEndpointURL.RawQuery = q.Encode()

	return authEndpointURL
}
