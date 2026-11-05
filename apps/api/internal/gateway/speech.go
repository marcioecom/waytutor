package gateway

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	SpeechModel  = "google/gemini-3.8-flash-tts"
	DefaultVoice = "Kore"

	speechInstructions = "Fale em português do Brasil como um professor didático: ritmo natural, " +
		"articulação clara, ênfase nos termos técnicos e pausas naturais entre as frases."
	outputFormat = "wav"
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

// SpeechCacheKey depends on everything that changes the audio, so a new model or
// a new instruction produces a new object instead of a stale hit.
func SpeechCacheKey(voice, text string) string {
	source := strings.Join([]string{"tts/v1", SpeechModel, voice, outputFormat, speechInstructions, text}, "|")
	sum := sha256.Sum256([]byte(source))
	return "tts/" + hex.EncodeToString(sum[:]) + ".wav"
}

func (c *Client) Speak(ctx context.Context, text, voice string) ([]byte, error) {
	body, err := c.post(ctx, "/v4/ai/speech-model", map[string]string{
		"text":         text,
		"voice":        voice,
		"outputFormat": outputFormat,
		"instructions": speechInstructions,
	})
	if err != nil {
		return nil, err
	}

	var payload struct {
		Audio string `json:"audio"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("%w: decode response: %s", ErrUpstream, err)
	}

	audio, err := base64.StdEncoding.DecodeString(payload.Audio)
	if err != nil || !isWav(audio) {
		return nil, fmt.Errorf("%w: response carries no wav audio", ErrUpstream)
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
