package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"gphoto-dler/cli/state"
	"gphoto-dler/google"
	"gphoto-dler/handler"

	"github.com/lmittmann/tint"
	"github.com/pkg/browser"
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
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("OK")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/start", handler.Start(googleClient))
	mux.HandleFunc("/callback", handler.Callback(googleClient))

	srv := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return
			}

			slog.Error(
				"failed to get media items",
				slog.String("error", err.Error()),
			)

			os.Exit(1)
		}
	}()
	defer srv.Close()

	for {
		httpClient := &http.Client{}
		resp, err := httpClient.Get("http://" + srv.Addr + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		resp.Body.Close()

		slog.Debug(
			"サーバの起動を待っています...",
			slog.String("addr", srv.Addr),
		)
		time.Sleep(10 * time.Millisecond)
	}

	slog.Info(
		"サーバが起動しました",
		slog.Group("endpoints",
			slog.String("start", fmt.Sprintf("http://%s/start", srv.Addr)),
			slog.String("callback", fmt.Sprintf("http://%s/callback", srv.Addr)),
		),
	)

	u := googleClient.BuildAuthURL()
	if err := browser.OpenURL(u.String()); err != nil {
		slog.Warn(
			"ブラウザを開くことができませんでした。以下のURLを開いてください。",
			slog.String("error", err.Error()),
		)
		fmt.Println(u.String())
	}

	// FIXME: チャンネル使ったほうがいいかも
	for {
		time.Sleep(1 * time.Second)

		if state.State.IsAuthenticated() {
			break
		}
	}

	slog.Info(
		"認証情報を更新するには、以下のURLから再認証してください。",
	)
	fmt.Println(u.String())
	for {
		time.Sleep(1 * time.Second)

		fmt.Print(state.State.StatusText() + "\033[5A")
	}

	medias, err := googleClient.GetMediaItems(&google.Token{
		TokenType:   "Bearer",
		AccessToken: state.State.AccessToken(),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(medias))
}

func findAvailablePort() string {
	for i := 49152; i < 65535; i++ {
		l, err := net.Listen("tcp", "127.0.0.1:"+fmt.Sprintf("%d", i))
		if err != nil {
			continue
		}

		defer l.Close()

		addr := l.Addr().String()

		return addr[len(addr)-4:]
	}

	panic("no available port")
}
