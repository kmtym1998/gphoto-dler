package handler

import (
	"log/slog"
	"net/http"

	"github.com/kmtym1998/gphoto-dler/google"
)

// Googleの認可エンドポイントにリダイレクトさせる
func Start(client *google.Client) func(w http.ResponseWriter, req *http.Request) {
	u := client.BuildAuthURL()

	// 認可エンドポイントにリダイレクト
	return func(w http.ResponseWriter, req *http.Request) {
		slog.Info(
			"redirecting to: " + u.String(),
		)

		http.Redirect(w, req, u.String(), http.StatusFound)
	}
}
