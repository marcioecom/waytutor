package gateway

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

func clientFor(t *testing.T, status int, body string) *Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	return New("key", server.URL)
}

func speakWith(t *testing.T, status int, body string) ([]byte, error) {
	t.Helper()
	return clientFor(t, status, body).Speak(context.Background(), "texto", DefaultVoice)
}

func TestSpeakDecodesTheAudio(t *testing.T) {
	body := `{"audio":"` + base64.StdEncoding.EncodeToString(wavBytes()) + `"}`

	audio, err := speakWith(t, http.StatusOK, body)
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
			if _, err := speakWith(t, tt.status, "{}"); !errors.Is(err, tt.want) {
				t.Fatalf("Speak() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestSpeakRejectsAudioThatIsNotWav(t *testing.T) {
	body := `{"audio":"` + base64.StdEncoding.EncodeToString([]byte("not a wav")) + `"}`

	if _, err := speakWith(t, http.StatusOK, body); !errors.Is(err, ErrUpstream) {
		t.Fatalf("Speak() error = %v, want ErrUpstream", err)
	}
}

func TestSpeakRejectsAnEmptyResponse(t *testing.T) {
	if _, err := speakWith(t, http.StatusOK, `{"audio":""}`); !errors.Is(err, ErrUpstream) {
		t.Fatalf("Speak() error = %v, want ErrUpstream", err)
	}
}

func TestSpeakSendsTheSpeechProtocolHeaders(t *testing.T) {
	var received http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = r.Header.Clone()
		_, _ = w.Write([]byte(`{"audio":"` + base64.StdEncoding.EncodeToString(wavBytes()) + `"}`))
	}))
	t.Cleanup(server.Close)

	if _, err := New("key", server.URL).Speak(context.Background(), "texto", DefaultVoice); err != nil {
		t.Fatalf("Speak() error = %v", err)
	}

	// The gateway rejects the request with 400 without these three.
	for name, want := range map[string]string{
		"ai-gateway-protocol-version":           "0.0.1",
		"ai-speech-model-specification-version": "4",
		"ai-model-id":                           SpeechModel,
	} {
		if got := received.Get(name); got != want {
			t.Errorf("header %s = %q, want %q", name, got, want)
		}
	}
}

func TestSpeechCacheKeyDependsOnVoiceAndText(t *testing.T) {
	base := SpeechCacheKey(DefaultVoice, "texto")
	if base != SpeechCacheKey(DefaultVoice, "texto") {
		t.Error("SpeechCacheKey is not stable for the same input")
	}
	if base == SpeechCacheKey("Puck", "texto") {
		t.Error("SpeechCacheKey ignores the voice")
	}
	if base == SpeechCacheKey(DefaultVoice, "outro texto") {
		t.Error("SpeechCacheKey ignores the text")
	}
	if base[len(base)-4:] != ".wav" {
		t.Errorf("SpeechCacheKey = %q, want a .wav key", base)
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

func TestSummarizeReturnsTheMessage(t *testing.T) {
	body := `{"choices":[{"message":{"content":"  Resumo do material.  "}}]}`

	summary, err := clientFor(t, http.StatusOK, body).Summarize(context.Background(), "material")
	if err != nil {
		t.Fatalf("Summarize() error = %v", err)
	}
	if summary != "Resumo do material." {
		t.Errorf("Summarize() = %q, want the trimmed message", summary)
	}
}

func TestSummarizeRejectsAnEmptyAnswer(t *testing.T) {
	for _, body := range []string{`{"choices":[]}`, `{"choices":[{"message":{"content":"   "}}]}`, `{}`} {
		if _, err := clientFor(t, http.StatusOK, body).Summarize(context.Background(), "material"); !errors.Is(err, ErrUpstream) {
			t.Fatalf("Summarize() error = %v for body %s, want ErrUpstream", err, body)
		}
	}
}

func TestSummarizeRejectsATruncatedAnswer(t *testing.T) {
	body := `{"choices":[{"message":{"content":"Primeiro, na seção de introdução e modelos de dados, vemos"},"finish_reason":"length"}]}`

	if _, err := clientFor(t, http.StatusOK, body).Summarize(context.Background(), "material"); !errors.Is(err, ErrUpstream) {
		t.Fatalf("Summarize() error = %v, want ErrUpstream", err)
	}
}
