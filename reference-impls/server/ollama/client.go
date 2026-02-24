package ollama

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	_ "embed"

	pb "github.com/dshearer/modelhawk/gen/go/v1"
	"github.com/ollama/ollama/api"
)

//go:embed prompt.md
var prompt string

func makePointer[T any](v T) *T {
	return &v
}

type modelResponse struct {
	Malicious bool   `json:"malicious"`
	Why       string `json:"why"`
}

var responseSchema = map[string]any{
	"$schema":  "http://json-schema.org/draft-07/schema#",
	"title":    "Malicious Detection Response",
	"type":     "object",
	"required": []string{"malicious"},
	"properties": map[string]any{
		"malicious": map[string]any{
			"type":        "boolean",
			"description": "Whether the tool call is malicious",
		},
		"why": map[string]any{
			"type":        "string",
			"description": "Explanation of why the tool call is malicious (required when malicious is true)",
		},
	},
	"if": map[string]any{
		"properties": map[string]any{
			"malicious": map[string]any{
				"const": true,
			},
		},
	},
	"then": map[string]any{
		"required": []string{"malicious", "why"},
	},
	"additionalProperties": false,
}

var responseSchemaSerialized = sync.OnceValue(func() []byte {
	s, err := json.Marshal(responseSchema)
	if err != nil {
		panic(err)
	}
	return s
})

// Client handles communication with a local Ollama server
type Client struct {
	msgLogger *msgLogger
	ollamaClient *api.Client
	conversation []api.Message
}

// NewClient creates a new Ollama client
func NewClient(base *url.URL, logDir string) (*Client, error) {
	// make system prompt message
	sysMsg := api.Message{
		Role:    "system",
		Content: prompt,
	}

	var ml *msgLogger
	if logDir != "" {
		var err error
		ml, err = newMsgLogger(logDir)
		if err != nil {
			return nil, err
		}
	}

	return &Client{
		msgLogger: ml,
		ollamaClient: api.NewClient(base, &http.Client{Timeout: 30 * time.Second}),
		conversation: []api.Message{sysMsg},
	}, nil
}

// CheckConnection verifies the connection to Ollama
func (c *Client) CheckConnection(ctx context.Context) error {
	if prompt == "" {
		panic("Prompt is empty")
	}
	if _, err := c.ollamaClient.List(ctx); err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	return nil
}

// Query sends a prompt to Ollama and returns the response
// If schema is not nil, it will be used to enforce structured output
func (c *Client) WantsToCallTool(ctx context.Context, req *pb.WantsToCallToolRequest) (*modelResponse, error) {
	// make messsage
	req.App = nil
	reqSerialized, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize WantsToCallToolRequest: %w", err)
	}
	newMsg := api.Message{
		Role:    "user",
		Content: string(reqSerialized),
	}
	c.conversation = append(c.conversation, newMsg)
	ollamaReq := api.ChatRequest{
		Model:    "llama3.2:latest",
		Format:   responseSchemaSerialized(),
		Messages: c.conversation,
		Stream:   makePointer(false),
	}

	// send it
	if c.msgLogger != nil {
		if err := c.msgLogger.LogRequest(&ollamaReq); err != nil {
			panic(err)
		}
	}
	var respMsg api.Message
	if err := c.ollamaClient.Chat(ctx, &ollamaReq, func(cr api.ChatResponse) error {
		if c.msgLogger != nil {
			if err := c.msgLogger.LogResponse(&cr); err != nil {
				panic(err)
			}
		}
		if !cr.Done {
			return errors.New("model did not finish making response")
		}
		if cr.Message.Role != "assistant" {
			return fmt.Errorf("unexpected role in response: %s", cr.Message.Role)
		}
		respMsg = cr.Message
		return nil
	}); err != nil {
		c.conversation = c.conversation[:len(c.conversation)]
		return nil, fmt.Errorf("failed to send request to Ollama: %w", err)
	}

	// save response to conversation
	c.conversation = append(c.conversation, respMsg)

	// parse it
	var parsedResp modelResponse
	if err := json.Unmarshal([]byte(respMsg.Content), &parsedResp); err != nil {
		return nil, fmt.Errorf("failed to parse response from model: %s\n\n%s", err, respMsg.Content)
	}

	return &parsedResp, nil
}
