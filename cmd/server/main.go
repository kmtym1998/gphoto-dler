package main

import (
	"log/slog"
	"net/http"
	"os"

	"gphoto-dler/google"
	"gphoto-dler/handler"

	"github.com/lmittmann/tint"
)

const (
	readPhotosScope = "https://www.googleapis.com/auth/photoslibrary.readonly"
	redirectURI     = "http://localhost:8080/callback"
)

func main() {
	slog.SetDefault(
		slog.New(
			tint.NewHandler(
				os.Stdout,
				&tint.Options{
					AddSource:  true,
					Level:      slog.LevelDebug,
					TimeFormat: "2006-01-02 15:04:05",
				},
			),
		),
	)

	googleClient, err := google.NewClient(
		readPhotosScope,
		redirectURI,
	)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/start", handler.Start(googleClient))
	http.HandleFunc("/callback", handler.Callback(googleClient))
	slog.Info(
		"starting server...",
		slog.Group("endpoints",
			slog.String("start", "http://localhost:8080/start"),
			slog.String("callback", "http://localhost:8080/callback"),
		),
	)

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		panic(err)
	}
}
