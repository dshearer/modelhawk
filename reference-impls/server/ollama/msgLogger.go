package ollama

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/ollama/ollama/api"
)

type msgLogger struct {
	dir     string
	counter int
	mu      sync.Mutex
}

// newMsgLogger creates a new client logger that writes to the specified directory.
// It scans for existing request_nnnn.json and response_nnnn.json files and starts
// counting from the highest number found.
func newMsgLogger(dir string) (*msgLogger, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Find the highest existing counter
	highestCounter, err := findHighestCounter(dir)
	if err != nil {
		return nil, err
	}

	return &msgLogger{
		dir:     dir,
		counter: highestCounter,
	}, nil
}

// findHighestCounter scans the directory for files matching the pattern
// and returns the highest counter value found, or 0 if none exist.
func findHighestCounter(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to read log directory: %w", err)
	}

	pattern := regexp.MustCompile(`^(request|response)_(\d+)\.json$`)
	maxCounter := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		matches := pattern.FindStringSubmatch(entry.Name())
		if len(matches) == 3 {
			counter, err := strconv.Atoi(matches[2])
			if err != nil {
				continue
			}
			if counter > maxCounter {
				maxCounter = counter
			}
		}
	}

	return maxCounter, nil
}

// LogRequest writes a request to a file with the next counter value.
func (cl *msgLogger) LogRequest(msg *api.ChatRequest) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	cl.counter++
	filename := fmt.Sprintf("request_%04d.json", cl.counter)
	return cl.writeJSON(filename, msg)
}

// logResponse writes a response to a file with the current counter value.
func (cl *msgLogger) LogResponse(msg *api.ChatResponse) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	filename := fmt.Sprintf("response_%04d.json", cl.counter)
	return cl.writeJSON(filename, msg)
}

// writeJSON marshals data to JSON and writes it to a file.
func (cl *msgLogger) writeJSON(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	path := filepath.Join(cl.dir, filename)
	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}

	return nil
}
