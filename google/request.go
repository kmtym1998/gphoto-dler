package google

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

// 認可コードを使ってトークンリクエストをエンドポイントに送る
func (c *Client) TokenRequest(code string) (*Token, error) {
	values := url.Values{}
	values.Add("client_id", c.secret.Installed.ClientID)
	values.Add("client_secret", c.secret.Installed.ClientSecret)
	values.Add("grant_type", c.option.GrantType)

	// 取得した認可コードをトークンのリクエストにセット
	values.Add("code", code)
	values.Add("redirect_uri", c.option.RedirectURI)

	// PKCE用パラメータ
	values.Add("code_verifier", verifier)

	req, err := http.NewRequest(
		http.MethodPost,
		c.option.TokenEndpoint,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	slog.Info(string(body))

	var data Token
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// 取得したトークンを利用してリソースにアクセス
func (c *Client) GetMediaItems(token *Token) ([]byte, error) {
	const endpoint = "https://photoslibrary.googleapis.com/v1/mediaItems"

	req, err := http.NewRequest(
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return nil, err
	}
	// 取得したアクセストークンをHeaderにセットしてリソースサーバにリクエストを送る
	req.Header.Set("Authorization", token.TokenType+" "+token.AccessToken)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get media items")
	}

	return body, nil
}
