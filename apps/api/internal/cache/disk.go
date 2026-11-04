package cache

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Disk struct {
	dir string
}

func NewDisk(dir string) (*Disk, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create cache directory: %w", err)
	}
	return &Disk{dir: dir}, nil
}

func (d *Disk) Get(_ context.Context, key string) ([]byte, error) {
	audio, err := os.ReadFile(d.path(key))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, ErrMiss
	}
	if err != nil {
		return nil, fmt.Errorf("read cached audio: %w", err)
	}
	return audio, nil
}

func (d *Disk) Put(_ context.Context, key string, audio []byte) error {
	target := d.path(key)
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create cache directory: %w", err)
	}

	temp, err := os.CreateTemp(filepath.Dir(target), "audio-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(temp.Name())

	if _, err := temp.Write(audio); err != nil {
		temp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Rename(temp.Name(), target); err != nil {
		return fmt.Errorf("move audio into place: %w", err)
	}
	return nil
}

func (d *Disk) path(key string) string {
	clean := filepath.Clean("/" + strings.ReplaceAll(key, "\\", "/"))
	return filepath.Join(d.dir, clean)
}
