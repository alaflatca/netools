package storage

import (
	"fmt"
	"io"
	"os"
)

type storageWriter interface {
	Write(io.Writer) error
}

func openWriter() (io.WriteCloser, error) {
	file, err := os.OpenFile(storagePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open file '%s': %w", storagePath, err)
	}

	return file, nil
}
