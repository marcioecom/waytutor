package cache

import (
	"context"
	"testing"
)

func TestDiskRoundTrip(t *testing.T) {
	cache, err := NewDisk(t.TempDir())
	if err != nil {
		t.Fatalf("NewDisk() error = %v", err)
	}
	ctx := context.Background()

	if _, err := cache.Get(ctx, "tts/abc.wav"); err != ErrMiss {
		t.Fatalf("Get() error = %v, want ErrMiss", err)
	}
	if err := cache.Put(ctx, "tts/abc.wav", []byte("audio")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	audio, err := cache.Get(ctx, "tts/abc.wav")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if string(audio) != "audio" {
		t.Errorf("Get() = %q, want audio", audio)
	}
}

func TestDiskKeepsKeysInsideTheCacheDirectory(t *testing.T) {
	cache, err := NewDisk(t.TempDir())
	if err != nil {
		t.Fatalf("NewDisk() error = %v", err)
	}

	if err := cache.Put(context.Background(), "../../escape.wav", []byte("audio")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	audio, err := cache.Get(context.Background(), "../../escape.wav")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if string(audio) != "audio" {
		t.Errorf("Get() = %q, want audio", audio)
	}
}
