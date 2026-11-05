package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	SummaryModel = "google/gemini-3.8-flash"

	summaryPrompt = `Você é um professor didático de curso superior. Resuma o MATERIAL para revisão em áudio, em português do Brasil.

Regras:
- Use apenas o MATERIAL. Não traga conhecimento externo, não complete lacunas e não invente.
- Se um ponto não estiver no material, diga "isso não aparece no material".
- Tamanho: 10 a 15% do material, cerca de 2 a 4 minutos falados.
- Prosa falada, sem markdown, sem listas e sem títulos.
- Use marcadores de oralidade: "primeiro", "em seguida", "por fim".
- Comece pelo conceito central e siga a ordem do material.
- Quando ajudar, cite a seção de origem entre parênteses.`
)

func (c *Client) Summarize(ctx context.Context, material string) (string, error) {
	body, err := c.post(ctx, "/v1/chat/completions", map[string]any{
		"model":       SummaryModel,
		"temperature": 0.2,
		"max_tokens":  1400,
		"messages": []map[string]string{
			{"role": "system", "content": summaryPrompt},
			{"role": "user", "content": "MATERIAL:\n\n" + material},
		},
	})
	if err != nil {
		return "", err
	}

	var payload struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("%w: decode response: %s", ErrUpstream, err)
	}
	if len(payload.Choices) == 0 {
		return "", fmt.Errorf("%w: response carries no summary", ErrUpstream)
	}

	summary := strings.TrimSpace(payload.Choices[0].Message.Content)
	if summary == "" {
		return "", fmt.Errorf("%w: response carries no summary", ErrUpstream)
	}
	return summary, nil
}
