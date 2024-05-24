package service

import "github.com/kmtym1998/gphoto-dler/google"

type Service struct {
	googleClient *google.Client
}

func New(googleClient *google.Client) *Service {
	return &Service{
		googleClient: googleClient,
	}
}
