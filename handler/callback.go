package handler

import (
	"gphoto-dler/google"
	"log"
	"net/http"
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

		// トークンレスポンスのjsonからトークンだけ抜き出しリソースにリクエストを送る
		body, err := client.GetMediaItems(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(body); err != nil {
			log.Println(err)
		}
	}
}
