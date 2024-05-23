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
	"gphoto-dler/service"

	"github.com/lmittmann/tint"
	"github.com/pkg/browser"
)

// MEMO: 並行処理させずに 6 min 程度で 70 枚 & 500MB くらいのダウンロードが完了する
// goroutine 同時数 20 で 5 分ぐらいで 500 枚 & 2GB くらいのダウンロードが完了する
// goroutine 同時数 100 で 1.5 分ぐらいで 500 枚 & 2GB くらいのダウンロードが完了する
// goroutine 同時数 200 で ↑ と同じくらいのペース

const (
	readPhotosScope      = "https://www.googleapis.com/auth/photoslibrary.readonly"
	tokenRefreshInterval = 55 * time.Minute
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

	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: gphoto-dler <destination>")
		os.Exit(1)
	}

	destDir := args[0]

	port := findAvailablePort()
	redirectURI := "http://localhost:" + port + "/callback"

	googleClient, err := google.NewClient(
		readPhotosScope,
		redirectURI,
	)
	if err != nil {
		panic(err)
	}
	s := service.New(googleClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("OK")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/start", handler.Start(googleClient))
	mux.HandleFunc("/callback", handler.Callback(s))

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
	go func() {
		for {
			lines := 0

			time.Sleep(1 * time.Second)

			fmt.Print(state.State.StatusText())
			lines += 5

			now := time.Now()
			if now.Add(60*time.Minute - tokenRefreshInterval).After(state.State.ExpiredAt()) {
				if state.State.RefreshToken() == "" {
					fmt.Println("アクセストークンが切れそうです。再認証してください。")
					lines++
				} else {
					if err := s.RefreshToken(state.State.RefreshToken()); err != nil {
						fmt.Println(err)
						lines++
					}
				}
			}

			fmt.Printf("\033[%dA", lines)
		}
	}()

	if err := s.DownloadMediaItems(destDir); err != nil {
		s.ExportError(destDir, err)
	}
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
