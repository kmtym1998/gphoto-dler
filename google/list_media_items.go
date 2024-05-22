package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// 取得したトークンを利用してリソースにアクセス
func (c *Client) ListMediaItems(token *Token, pageSize int, nextPageToken string) (*ListMediaItemsResult, error) {
	// NOTE: https://developers.google.com/photos/library/reference/rest/v1/mediaItems/list?hl=ja
	const endpoint = "https://photoslibrary.googleapis.com/v1/mediaItems"

	u, _ := url.Parse(endpoint)
	q := u.Query()
	q.Add("pageSize", fmt.Sprint(pageSize))
	if nextPageToken != "" {
		q.Add("pageToken", nextPageToken)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
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

	var data ListMediaItemsResult
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
