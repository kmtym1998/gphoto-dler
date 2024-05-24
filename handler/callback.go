package handler

import (
	"net/http"

	"github.com/kmtym1998/gphoto-dler/service"
)

// 認可してからcallbackするところ
func Callback(s *service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "code not found", http.StatusBadRequest)
			return
		}

		if err := s.GetAndSaveNewToken(code); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write([]byte("<html><body><h1>認証完了</h1></body></html>")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
