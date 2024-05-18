package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	responseType = "code"
	redirectURI  = "http://localhost:8080/callback"
	grantType    = "authorization_code"

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

// Googleの認可エンドポイントにリダイレクトさせる
func start(w http.ResponseWriter, req *http.Request) {
	authEndpoint := oauth.authEndpoint

	values := url.Values{}
	values.Add("response_type", responseType)
	values.Add("client_id", oauth.clientID)
	values.Add("state", oauth.state)
	values.Add("scope", oauth.scope)
	values.Add("redirect_uri", redirectURI)

	// PKCE用パラメータ
	values.Add("code_challenge_method", oauth.codeChallengeMethod)
	values.Add("code_challenge", oauth.codeChallenge)

	// 認可エンドポイントにリダイレクト
	http.Redirect(w, req, authEndpoint+values.Encode(), 302)
}

// 認可してからcallbackするところ
func callback(w http.ResponseWriter, req *http.Request) {
	//クエリを取得
	query := req.URL.Query()

	// トークンをリクエストする
	result, err := tokenRequest(query)
	if err != nil {
		log.Println(err)
	}
	// トークンレスポンスのjsonからトークンだけ抜き出しリソースにリクエストを送る
	body, err := apiRequest(req, result["access_token"].(string))
	if err != nil {
		log.Println(err)
	}
	if _, err := w.Write(body); err != nil {
		log.Println(err)
	}
}

// 認可コードを使ってトークンリクエストをエンドポイントに送る
func tokenRequest(query url.Values) (map[string]interface{}, error) {
	tokenEndpoint := oauth.tokenEndpoint
	values := url.Values{}
	values.Add("client_id", oauth.clientID)
	values.Add("client_secret", oauth.clientSecret)
	values.Add("grant_type", grantType)

	// 取得した認可コードをトークンのリクエストにセット
	values.Add("code", query.Get("code"))
	values.Add("redirect_uri", redirectURI)

	// PKCE用パラメータ
	values.Add("code_verifier", verifier)

	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("request err: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("token response : %s", string(body))
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// 取得したトークンを利用してリソースにアクセス
func apiRequest(_ *http.Request, token string) ([]byte, error) {
	photoAPI := "https://photoslibrary.googleapis.com/v1/mediaItems"

	req, err := http.NewRequest("GET", photoAPI, nil)
	if err != nil {
		return nil, err
	}
	// 取得したアクセストークンをHeaderにセットしてリソースサーバにリクエストを送る
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("http status code is %d, err: %s", resp.StatusCode, err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return body, nil
}

func main() {
	setUp()
	http.HandleFunc("/start", start)
	http.HandleFunc("/callback", callback)
	log.Println("starting server...")
	log.Println("http://localhost:8080/start にアクセスしてください")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
