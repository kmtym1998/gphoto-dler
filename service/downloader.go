package service

import (
	"errors"
	"fmt"
	"gphoto-dler/cli/state"
	"gphoto-dler/google"
	"io"
	"net/http"
	"os"
	"strings"
)

func (s *Service) DownloadMediaItems() error {
	if state.State.AccessToken() == "" {
		return errors.New("access token is empty")
	}

	var errList []error

	result, err := s.googleClient.ListMediaItems(&google.Token{
		AccessToken: state.State.AccessToken(),
		TokenType:   "Bearer",
	}, 100, "")
	if err != nil {
		return err
	}
	nextPageToken := result.NextPageToken

	if err := batchDownloadMediaItems("downloads", result); err != nil {
		errList = append(errList, err)
	}

	for {
		if nextPageToken == "" {
			break
		}

		result, err = s.googleClient.ListMediaItems(&google.Token{
			AccessToken: state.State.AccessToken(),
			TokenType:   "Bearer",
		}, 100, result.NextPageToken)
		if err != nil {
			errList = append(errList, err)
		}
		nextPageToken = result.NextPageToken

		if err := batchDownloadMediaItems("downloads", result); err != nil {
			errList = append(errList, err)
		}
	}

	return errors.Join(errList...)
}

func batchDownloadMediaItems(dstDir string, list *google.ListMediaItemsResult) error {
	var errList []error

	for _, item := range list.MediaItems {
		var query string
		switch determineMediaType(item.MimeType) {
		case mediaTypeImage:
			query = "=d-w" + item.MediaMetadata.Width + "-h" + item.MediaMetadata.Height
		case mediaTypeVideo:
			if item.MediaMetadata.Video.Status == "READY" {
				query = "=dv-d"
			} else {
				query = "=d-w" + item.MediaMetadata.Width + "-h" + item.MediaMetadata.Height + "-no"
			}
		case mediaTypeUnknown:
			query = "=d-w" + item.MediaMetadata.Width + "-h" + item.MediaMetadata.Height
		}

		endpoint := item.BaseURL + query

		resp, err := http.Get(endpoint)
		if err != nil {
			errList = append(errList, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			errList = append(errList, fmt.Errorf("failed to download %s. status: %s", item.Filename, resp.Status))
			continue
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			errList = append(errList, err)
			continue
		}
		resp.Body.Close()

		if err := os.WriteFile(dstDir+"/"+item.Filename, b, os.ModePerm); err != nil {
			errList = append(errList, err)
			continue
		}
	}

	if len(errList) > 0 {
		return errors.Join(errList...)
	}

	return nil
}

type mediaType string

const (
	mediaTypeImage   mediaType = "image"
	mediaTypeVideo   mediaType = "video"
	mediaTypeUnknown mediaType = "unknown"
)

func determineMediaType(mimeType string) mediaType {
	if strings.HasPrefix(mimeType, "image/") {
		return mediaTypeImage
	}

	if strings.HasPrefix(mimeType, "video/") {
		return mediaTypeVideo
	}

	return mediaTypeUnknown
}
