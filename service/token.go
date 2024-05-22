package service

import (
	"errors"
	"gphoto-dler/cli/state"
	"time"
)

func (s *Service) RefreshToken(refreshToken string) error {
	if refreshToken == "" {
		return errors.New("refresh token is empty")
	}

	token, err := s.googleClient.TokenRequestByRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	state.State.SetAccessToken(token.AccessToken)
	state.State.SetRefreshToken("")
	state.State.SetExpiredAt(time.Duration(token.ExpiresIn) * time.Second)

	return nil
}

func (s *Service) GetAndSaveNewToken(authorizationCode string) error {
	if authorizationCode == "" {
		return errors.New("authorization code is empty")
	}

	// トークンをリクエストする
	token, err := s.googleClient.TokenRequestByAuthorizationCode(authorizationCode)
	if err != nil {
		return err
	}

	state.State.SetAccessToken(token.AccessToken)
	state.State.SetRefreshToken(token.RefreshToken)
	state.State.SetExpiredAt(time.Duration(token.ExpiresIn) * time.Second)

	return nil
}
