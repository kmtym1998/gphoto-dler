package handler

import (
	"net/http"

	"gphoto-dler/google"
)

// Googleの認可エンドポイントにリダイレクトさせる
func Start(client *google.Client) func(w http.ResponseWriter, req *http.Request) {
	u := client.BuildAuthURL()

	// 認可エンドポイントにリダイレクト
	return func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, u.String(), http.StatusFound)
	}
}
