package handler

import (
	"net/http"
	"time"

	"gphoto-dler/cli/state"
	"gphoto-dler/google"
)

// 認可してからcallbackするところ
func Callback(client *google.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "code not found", http.StatusBadRequest)
			return
		}

		// トークンをリクエストする
		token, err := client.TokenRequest(code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		state.State.SetAccessToken(token.AccessToken)
		state.State.SetRefreshToken(token.RefreshToken)
		state.State.SetExpiredAt(time.Duration(token.ExpiresIn) * time.Second)

		if _, err := w.Write([]byte("<html><body><h1>認証完了</h1></body></html>")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
