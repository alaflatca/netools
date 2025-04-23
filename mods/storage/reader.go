package storage

import (
	"fmt"
	"io"
	"os"
)

type storageReader interface {
	Read(io.Reader) error
}

func openReader() (io.ReadCloser, error) {
	file, err := os.Open(storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file  %w", err)
	}

	return file, nil
}
