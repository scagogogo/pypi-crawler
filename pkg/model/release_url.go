package model

import (
	"context"
	"fmt"
	"github.com/crawler-go-go-go/go-requests"
	"time"
)

// ReleaseURL 发布的URL
type ReleaseURL struct {
	CommentText       string    `json:"comment_text"`
	Digests           *Digests  `json:"digests"`
	Downloads         int       `json:"downloads"`
	Filename          string    `json:"filename"`
	HasSig            bool      `json:"has_sig"`
	Md5Digest         string    `json:"md5_digest"`
	Packagetype       string    `json:"packagetype"`
	PythonVersion     string    `json:"python_version"`
	RequiresPython    string    `json:"requires_python"`
	Size              int       `json:"size"`
	UploadTime        string    `json:"upload_time"`
	UploadTimeIso8601 time.Time `json:"upload_time_iso_8601"`
	URL               string    `json:"url"`
	Yanked            bool      `json:"yanked"`

	// TODO
	YankedReason any `json:"yanked_reason"`
}

func (x *ReleaseURL) Download(ctx context.Context) ([]byte, error) {
	if x.URL == "" {
		return nil, fmt.Errorf("do not has release file url")
	}
	return requests.GetBytes(ctx, x.URL)
}

type Digests struct {
	Blake2b_256 string `json:"blake2b_256"`
	MD5         string `json:"md5"`
	Sha256      string `json:"sha256"`
}
