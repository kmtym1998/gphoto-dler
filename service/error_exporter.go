package service

import (
	"fmt"
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

	s.writeError(f, err)
}

func (s *Service) writeError(f *os.File, err error) {
	if err == nil {
		return
	}

	spreadableErr, ok := err.(interface {
		Unwrap() []error
	})
	if !ok {
		fmt.Fprintln(f, err)
	}

	for _, e := range spreadableErr.Unwrap() {
		s.writeError(f, e)
	}
}
