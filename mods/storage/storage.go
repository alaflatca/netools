package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// netools
// 0, SSH 정보
// 1, vpn
// 2, reverse proxy
// 3, ping
// 4. termios logging

func Read(reader storageReader) error {
	file, err := openReader()
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("error closing file %v\n", closeErr)
		}
	}()

	return reader.Read(file)
}

func Write(writer storageWriter) error {
	file, err := openWriter()
	if err != nil {
		return fmt.Errorf("failed to open storage file: %w", err)
	}

	var closeErr error
	defer func() {
		if closeErr = file.Close(); closeErr != nil {
			log.Printf("error closing file: %v\n", closeErr)
		}
	}()

	if err = writer.Write(file); err != nil {
		return fmt.Errorf("failed to write to storage file: %w", err)
	}

	return closeErr
}

var storagePath string

func CreateStorage() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	storagePath = filepath.Join(homeDir, ".config")
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		log.Printf("create config path: '%s'\n", storagePath)
		storagePath = filepath.Join(storagePath, "netools")
		if err := os.MkdirAll(storagePath, 0700); err != nil {
			return err
		}
	} else {
		storagePath = filepath.Join(storagePath, "netools")
	}

	storageName := "netools.db"
	storagePath = filepath.Join(storagePath, storageName)
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		log.Printf("create storage file: '%s'\n", storagePath)
		f, err := os.OpenFile(storagePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0700)
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}
