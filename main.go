package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"gphoto-dler/google"
	"gphoto-dler/handler"

	"github.com/lmittmann/tint"
)

const (
	responseType    = "code"
	readPhotosScope = "https://www.googleapis.com/auth/photoslibrary.readonly"
	redirectURI     = "http://localhost:8080/callback"
	grantType       = "authorization_code"

	// https://tex2e.github.io/rfc-translater/html/rfc7636.html
	// 付録B. S256 code_challenge_methodの例 "
	verifier = "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
)

var secrets map[string]interface{}

var oauth struct {
	clientID            string
	clientSecret        string
	scope               string
	state               string
	codeChallengeMethod string
	codeChallenge       string
	authEndpoint        string
	tokenEndpoint       string
}

func readJson() {
	data, err := os.ReadFile("secrets/secret.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &secrets); err != nil {
		panic(err)
	}
}

func setUp() {
	readJson()

	oauth.clientID = secrets["installed"].(map[string]interface{})["client_id"].(string)
	oauth.clientSecret = secrets["installed"].(map[string]interface{})["client_secret"].(string)
	oauth.authEndpoint = "https://accounts.google.com/o/oauth2/v2/auth?"
	oauth.tokenEndpoint = "https://www.googleapis.com/oauth2/v4/token"
	oauth.state = "xyz"
	oauth.scope = "https://www.googleapis.com/auth/photoslibrary.readonly"
	oauth.codeChallengeMethod = "S256"

	// PKCE用に"dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"をSHA256+Base64URLエンコードしたものをセット
	oauth.codeChallenge = base64URLEncode()

}

// https://auth0.com/docs/authorization/flows/call-your-api-using-the-authorization-code-flow-with-pkce#javascript-sample
func base64URLEncode() string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

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

	setUp()
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
