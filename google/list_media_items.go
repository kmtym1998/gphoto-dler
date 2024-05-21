package google

import (
	"errors"
	"io"
	"log"
	"net/http"
)

// 取得したトークンを利用してリソースにアクセス
func (c *Client) GetMediaItems(token *Token) ([]byte, error) {
	// NOTE: https://developers.google.com/photos/library/reference/rest/v1/mediaItems/list?hl=ja
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
