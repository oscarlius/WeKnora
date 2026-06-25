package vlm

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Tencent/WeKnora/internal/models/provider"
	openai "github.com/sashabaranov/go-openai"
)

func TestVLMImageDetail(t *testing.T) {
	t.Run("MiniMax omits detail for compatibility", func(t *testing.T) {
		t.Setenv("VLM_IMAGE_DETAIL", "")
		got, send := vlmImageDetail(provider.ProviderMiniMax)
		if send {
			t.Fatalf("expected MiniMax detail to be omitted, got %q", got)
		}
	})

	t.Run("OpenAI keeps auto detail", func(t *testing.T) {
		t.Setenv("VLM_IMAGE_DETAIL", "")
		got, send := vlmImageDetail(provider.ProviderOpenAI)
		if !send || got != openai.ImageURLDetailAuto {
			t.Fatalf("expected OpenAI detail auto to be sent, got %q send=%v", got, send)
		}
	})

	t.Run("generic OpenAI-compatible omits detail", func(t *testing.T) {
		t.Setenv("VLM_IMAGE_DETAIL", "")
		got, send := vlmImageDetail(provider.ProviderGeneric)
		if send {
			t.Fatalf("expected generic detail to be omitted, got %q", got)
		}
	})

	t.Run("environment override wins", func(t *testing.T) {
		t.Setenv("VLM_IMAGE_DETAIL", "low")
		got, send := vlmImageDetail(provider.ProviderMiniMax)
		if !send || got != openai.ImageURLDetailLow {
			t.Fatalf("expected env override low to be sent, got %q send=%v", got, send)
		}
	})

	t.Run("environment can force omit", func(t *testing.T) {
		t.Setenv("VLM_IMAGE_DETAIL", "omit")
		got, send := vlmImageDetail(provider.ProviderOpenAI)
		if send {
			t.Fatalf("expected env override omit, got %q", got)
		}
	})
}

func TestVLMImageDetailOmittedFromJSON(t *testing.T) {
	part := openai.ChatMessagePart{
		Type: openai.ChatMessagePartTypeImageURL,
		ImageURL: &openai.ChatMessageImageURL{
			URL: "data:image/png;base64,AA==",
		},
	}
	body, err := json.Marshal(part)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "detail") {
		t.Fatalf("expected image detail to be omitted, got %s", string(body))
	}
}
