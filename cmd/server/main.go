package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"gphoto-dler/google"
	"gphoto-dler/handler"

	"github.com/lmittmann/tint"
)

const (
	readPhotosScope = "https://www.googleapis.com/auth/photoslibrary.readonly"
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

	port := findAvailablePort()
	redirectURI := "http://localhost:" + port + "/callback"

	googleClient, err := google.NewClient(
		readPhotosScope,
		redirectURI,
	)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/start", handler.Start(googleClient))
	mux.HandleFunc("/callback", handler.Callback(googleClient))

	srv := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	defer srv.Close()

	slog.Info(
		"starting server...",
		slog.Group("endpoints",
			slog.String("start", fmt.Sprintf("http://%s/start", srv.Addr)),
			slog.String("callback", fmt.Sprintf("http://%s/callback", srv.Addr)),
		),
	)
}

func findAvailablePort() string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	addr := l.Addr().String()

	return addr[len(addr)-4:]
}
