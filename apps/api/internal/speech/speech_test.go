package speech

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func wavBytes() []byte {
	audio := make([]byte, 44)
	copy(audio[0:4], "RIFF")
	copy(audio[8:12], "WAVE")
	audio[4] = 36
	return audio
}

func synthesizeWith(t *testing.T, status int, body string) ([]byte, error) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	return New("key", server.URL).Speak(context.Background(), "texto", DefaultVoice)
}

func TestSpeakDecodesTheAudio(t *testing.T) {
	body := `{"audio":"` + base64.StdEncoding.EncodeToString(wavBytes()) + `"}`

	audio, err := synthesizeWith(t, http.StatusOK, body)
	if err != nil {
		t.Fatalf("Speak() error = %v", err)
	}
	if len(audio) != len(wavBytes()) {
		t.Errorf("len(audio) = %d, want %d", len(audio), len(wavBytes()))
	}
}

func TestSpeakMapsProviderFailures(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   error
	}{
		{"credentials", http.StatusUnauthorized, ErrAuth},
		{"forbidden", http.StatusForbidden, ErrAuth},
		{"payment", http.StatusPaymentRequired, ErrQuota},
		{"rate limit", http.StatusTooManyRequests, ErrQuota},
		{"server error", http.StatusInternalServerError, ErrUpstream},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := synthesizeWith(t, tt.status, "{}"); !errors.Is(err, tt.want) {
				t.Fatalf("Speak() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestSpeakRejectsAudioThatIsNotWav(t *testing.T) {
	body := `{"audio":"` + base64.StdEncoding.EncodeToString([]byte("not a wav")) + `"}`

	if _, err := synthesizeWith(t, http.StatusOK, body); !errors.Is(err, ErrUpstream) {
		t.Fatalf("Speak() error = %v, want ErrUpstream", err)
	}
}

func TestSpeakRejectsAnEmptyResponse(t *testing.T) {
	if _, err := synthesizeWith(t, http.StatusOK, `{"audio":""}`); !errors.Is(err, ErrUpstream) {
		t.Fatalf("Speak() error = %v, want ErrUpstream", err)
	}
}

func TestCacheKeyDependsOnVoiceAndText(t *testing.T) {
	base := CacheKey(DefaultVoice, "texto")
	if base != CacheKey(DefaultVoice, "texto") {
		t.Error("CacheKey is not stable for the same input")
	}
	if base == CacheKey("Puck", "texto") {
		t.Error("CacheKey ignores the voice")
	}
	if base == CacheKey(DefaultVoice, "outro texto") {
		t.Error("CacheKey ignores the text")
	}
	if base[len(base)-4:] != ".wav" {
		t.Errorf("CacheKey = %q, want a .wav key", base)
	}
}

func TestValidVoice(t *testing.T) {
	if !ValidVoice(DefaultVoice) {
		t.Errorf("ValidVoice(%q) = false", DefaultVoice)
	}
	if ValidVoice("Nenhuma") {
		t.Error("ValidVoice accepted an unknown voice")
	}
}
