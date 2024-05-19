package state

import (
	"sync"
	"time"
)

type GlobalState struct {
	accessToken  string
	refreshToken string
	// アクセストークンが無効になる時間
	expiredAt time.Time

	mutex sync.Mutex
}

var State GlobalState = GlobalState{
	mutex: sync.Mutex{},
}

var StateChannel chan GlobalState = make(chan GlobalState)

// リフレッシュが必要かどうか (アクセストークンが無効になる5分前にリフレッシュ)
func (s *GlobalState) ShouldRefresh() bool {
	now := time.Now()

	return now.After(s.expiredAt.Add(-5 * time.Minute))
}

// 再認証が必要かどうか (アクセストークンが切れそう && リフレッシュトークンも無い)
func (s *GlobalState) ShouldRotate() bool {
	return s.ShouldRefresh() && s.refreshToken == ""
}

// 認証済みかどうか
func (s *GlobalState) IsAuthenticated() bool {
	return s.accessToken != ""
}

// accessTokenをセット
func (s *GlobalState) SetAccessToken(token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.accessToken = token
}

// accessTokenを取得
func (s *GlobalState) AccessToken() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.accessToken
}

// refreshTokenをセット
func (s *GlobalState) SetRefreshToken(token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.refreshToken = token
}

// refreshTokenを取得
func (s *GlobalState) RefreshToken() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.refreshToken
}

// expiredAtをセット
func (s *GlobalState) SetExpiredAt(expiredAfter time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.expiredAt = time.Now().Add(expiredAfter)
}

// expiredAtを取得
func (s *GlobalState) ExpiredAt() time.Time {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.expiredAt
}
