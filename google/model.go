package google

import "time"

type ListMediaItemsResult struct {
	MediaItems    []MediaItems `json:"mediaItems,omitempty"`
	NextPageToken string       `json:"nextPageToken,omitempty"`
}

type MediaItems struct {
	ID            string        `json:"id,omitempty"`
	ProductURL    string        `json:"productUrl,omitempty"`
	BaseURL       string        `json:"baseUrl,omitempty"`
	MimeType      string        `json:"mimeType,omitempty"`
	MediaMetadata MediaMetadata `json:"mediaMetadata,omitempty"`
	Filename      string        `json:"filename,omitempty"`
}

type MediaMetadata struct {
	CreationTime time.Time `json:"creationTime,omitempty"`
	Width        string    `json:"width,omitempty"`
	Height       string    `json:"height,omitempty"`
	// 写真のメディアタイプのメタデータ
	Photo Photo `json:"photo,omitempty"`
	// 動画のメディアタイプのメタデータ
	Video Video `json:"video,omitempty"`
}
type Video struct {
	CameraMake  string  `json:"cameraMake,omitempty"`
	CameraModel string  `json:"cameraModel,omitempty"`
	FPS         float64 `json:"fps,omitempty"`
	Status      string  `json:"status,omitempty"`
}

type Photo struct {
	CameraMake      string  `json:"cameraMake,omitempty"`
	CameraModel     string  `json:"cameraModel,omitempty"`
	FocalLength     float64 `json:"focalLength,omitempty"`
	ApertureFNumber float64 `json:"apertureFNumber,omitempty"`
	ISOEquivalent   int     `json:"isoEquivalent,omitempty"`
	ExposureTime    string  `json:"exposureTime,omitempty"`
}
