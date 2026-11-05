package speech

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	Model          = "google/gemini-3.8-flash-tts"
	APIBaseURL     = "https://ai-gateway.vercel.sh"
	DefaultVoice   = "Kore"
	RequestTimeout = 120 * time.Second

	instructions = "Fale em português do Brasil como um professor didático: ritmo natural, " +
		"articulação clara, ênfase nos termos técnicos e pausas naturais entre as frases."
	outputFormat    = "wav"
	maxResponseBody = 32 << 20
)

var (
	ErrAuth     = errors.New("speech provider rejected the credentials")
	ErrQuota    = errors.New("speech provider is out of quota")
	ErrUpstream = errors.New("speech provider failed")
)

var voices = []string{
	"Kore", "Puck", "Aoede", "Charon", "Fenrir", "Zephyr", "Leda", "Orus",
	"Callirrhoe", "Autonoe", "Enceladus", "Iapetus", "Umbriel", "Algieba",
	"Despina", "Erinome", "Algenib", "Rasalgethi", "Laomedeia", "Achernar",
	"Alnilam", "Schedar", "Gacrux", "Pulcherrima", "Achird", "Zubenelgenubi",
	"Vindemiatrix", "Sadachbia", "Sadaltager", "Sulafat",
}

func ValidVoice(voice string) bool {
	for _, known := range voices {
		if known == voice {
			return true
		}
	}
	return false
}

func Voices() []string {
	return append([]string(nil), voices...)
}

type Synthesizer struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func New(apiKey, baseURL string) *Synthesizer {
	return &Synthesizer{
		apiKey:  apiKey,
		baseURL: strings.TrimSuffix(baseURL, "/"),
		client:  &http.Client{Timeout: RequestTimeout},
	}
}

// CacheKey depends on everything that changes the audio, so a new model or a
// new instruction produces a new object instead of a stale hit.
func CacheKey(voice, text string) string {
	source := strings.Join([]string{"tts/v1", Model, voice, outputFormat, instructions, text}, "|")
	sum := sha256.Sum256([]byte(source))
	return "tts/" + hex.EncodeToString(sum[:]) + ".wav"
}

func (s *Synthesizer) Speak(ctx context.Context, text, voice string) ([]byte, error) {
	body, err := json.Marshal(map[string]string{
		"text":         text,
		"voice":        voice,
		"outputFormat": outputFormat,
		"instructions": instructions,
	})
	if err != nil {
		return nil, fmt.Errorf("encode speech request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/v4/ai/speech-model", strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("build speech request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("ai-gateway-protocol-version", "0.0.1")
	req.Header.Set("ai-speech-model-specification-version", "4")
	req.Header.Set("ai-model-id", Model)

	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUpstream, err)
	}
	defer res.Body.Close()

	switch {
	case res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden:
		return nil, fmt.Errorf("%w: %s", ErrAuth, res.Status)
	case res.StatusCode == http.StatusPaymentRequired || res.StatusCode == http.StatusTooManyRequests:
		return nil, fmt.Errorf("%w: %s", ErrQuota, res.Status)
	case res.StatusCode != http.StatusOK:
		detail, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return nil, fmt.Errorf("%w: %s: %s", ErrUpstream, res.Status, detail)
	}

	var payload struct {
		Audio    string `json:"audio"`
		Warnings []struct {
			Details string `json:"details"`
		} `json:"warnings"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, maxResponseBody)).Decode(&payload); err != nil {
		return nil, fmt.Errorf("%w: decode response: %s", ErrUpstream, err)
	}

	audio, err := base64.StdEncoding.DecodeString(payload.Audio)
	if err != nil || len(audio) == 0 {
		return nil, fmt.Errorf("%w: response carries no audio", ErrUpstream)
	}
	if !isWav(audio) {
		return nil, fmt.Errorf("%w: response is not a wav file", ErrUpstream)
	}

	return audio, nil
}

func isWav(audio []byte) bool {
	if len(audio) < 12 {
		return false
	}
	if string(audio[0:4]) != "RIFF" || string(audio[8:12]) != "WAVE" {
		return false
	}
	return binary.LittleEndian.Uint32(audio[4:8]) > 0
}
