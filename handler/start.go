package handler

import (
	"gphoto-dler/google"
	"net/http"
	"net/url"
)

// Googleの認可エンドポイントにリダイレクトさせる
func Start(client *google.Client) func(w http.ResponseWriter, req *http.Request) {
	opt := client.Option()
	secret := client.Secret()

	authEndpointURL, err := url.Parse(opt.AuthEndpoint)
	if err != nil {
		panic(err)
	}

	q := authEndpointURL.Query()
	q.Add("client_id", secret.Installed.ClientID)
	q.Add("state", opt.State)
	q.Add("redirect_uri", opt.RedirectURI)
	q.Add("response_type", opt.ResponseType)
	q.Add("scope", opt.Scope)

	// PKCE用パラメータ
	q.Add("code_challenge", opt.CodeChallenge)
	q.Add("code_challenge_method", opt.CodeChallengeMethod)

	authEndpointURL.RawQuery = q.Encode()

	// 認可エンドポイントにリダイレクト
	return func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, authEndpointURL.String(), http.StatusFound)
	}
}
