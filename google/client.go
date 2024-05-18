package google

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
)

const (
	verifier = "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
)

type Client struct {
	secret Secret
	option ClientOption
}

type Secret struct {
	Installed Installed `json:"installed"`
}

type Installed struct {
	ClientID                string   `json:"client_id"`
	ProjectID               string   `json:"project_id"`
	AuthURI                 string   `json:"auth_uri"`
	TokenURI                string   `json:"token_uri"`
	AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
	ClientSecret            string   `json:"client_secret"`
	RedirectUris            []string `json:"redirect_uris"`
}

type ClientOption struct {
	AuthEndpoint        string
	CodeChallenge       string
	CodeChallengeMethod string
	GrantType           string
	RedirectURI         string
	ResponseType        string
	Scope               string
	State               string
	TokenEndpoint       string
	Verifier            string
}

func NewClient(scope, redirectURI string) (*Client, error) {
	b, err := os.ReadFile("secrets/secret.json")
	if err != nil {
		return nil, err
	}

	var secret Secret
	if err := json.Unmarshal(b, &secret); err != nil {
		return nil, err
	}

	o := ClientOption{
		AuthEndpoint:        "https://accounts.google.com/o/oauth2/v2/auth",
		CodeChallenge:       base64URLEncode(),
		CodeChallengeMethod: "S256",
		GrantType:           "authorization_code",
		RedirectURI:         redirectURI,
		ResponseType:        "code",
		Scope:               scope,
		State:               "xyz",
		TokenEndpoint:       "https://www.googleapis.com/oauth2/v4/token",
		// https://tex2e.github.io/rfc-translater/html/rfc7636.html
		// 付録B. S256 code_challenge_methodの例 "
		Verifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
	}

	return &Client{
		secret: secret,
		option: o,
	}, nil
}

func (c *Client) Option() ClientOption {
	return c.option
}

func (c *Client) Secret() Secret {
	return c.secret
}

// https://auth0.com/docs/authorization/flows/call-your-api-using-the-authorization-code-flow-with-pkce#javascript-sample
func base64URLEncode() string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
