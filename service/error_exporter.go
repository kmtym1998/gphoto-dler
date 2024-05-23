package service

import (
	"encoding/json"
	"fmt"
	"gphoto-dler/cli/state"
	"log/slog"
	"os"
	"time"
)

func (s *Service) ExportError(destDir string, err error) {
	if err == nil {
		return
	}

	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		slog.Error("ディレクトリを作成できませんでした", slog.String("error", err.Error()))
		return
	}

	now := time.Now()
	f, err := os.Create(fmt.Sprintf(
		"%s/%d_error.log",
		destDir,
		now.Unix(),
	))
	if err != nil {
		slog.Error("ファイルを作成できませんでした", slog.String("error", err.Error()))
		return
	}

	items := state.State.FailedItems()
	b, err := json.Marshal(items)
	if err != nil {
		slog.Error("JSONエンコードに失敗しました", slog.String("error", err.Error()))
		return
	}

	if _, err := f.Write(b); err != nil {
		slog.Error("ファイルに書き込めませんでした", slog.String("error", err.Error()))
		return
	}
}
