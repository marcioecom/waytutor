package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	SummaryModel = "openai/gpt-6-luna"

	// The model reasons before answering, and reasoning tokens count against the
	// budget. A summary of 2 to 4 spoken minutes costs ~880 tokens.
	summaryMaxTokens = 2000

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
		"max_tokens":  summaryMaxTokens,
		"messages": []map[string]string{
			{"role": "system", "content": summaryPrompt},
			{"role": "user", "content": "MATERIAL:\n\n" + material},
		},
	}, nil)
	if err != nil {
		return "", err
	}

	var payload struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("%w: decode response: %s", ErrUpstream, err)
	}
	if len(payload.Choices) == 0 {
		return "", fmt.Errorf("%w: response carries no summary", ErrUpstream)
	}

	choice := payload.Choices[0]
	// A cut summary reads like a broken one, so it is a failure, not a result.
	if strings.EqualFold(choice.FinishReason, "length") {
		return "", fmt.Errorf("%w: summary hit the token budget", ErrUpstream)
	}

	summary := strings.TrimSpace(choice.Message.Content)
	if summary == "" {
		return "", fmt.Errorf("%w: response carries no summary", ErrUpstream)
	}
	return summary, nil
}
