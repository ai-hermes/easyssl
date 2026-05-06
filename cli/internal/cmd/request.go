package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"easyssl/cli/internal/client"
)

func runRequest(method, reqPath string, query map[string]string, body any, auth client.AuthMode) error {
	if err := requireAuth(); err != nil {
		return err
	}
	c, err := newAPIClient()
	if err != nil {
		return err
	}
	res, err := c.Do(method, reqPath, query, body, auth)
	if err != nil {
		return parseAPIError(err)
	}

	var data any
	if len(res.Envelope.Data) > 0 {
		if err := json.Unmarshal(res.Envelope.Data, &data); err != nil {
			data = string(res.Envelope.Data)
		}
	}
	if data == nil {
		data = map[string]any{}
	}
	if err := printOutput(map[string]any{
		"code": res.Envelope.Code,
		"msg":  res.Envelope.Msg,
		"data": data,
	}); err != nil {
		return exitErr(5, err)
	}
	return nil
}

func readJSONInput(filePath, inline string) (map[string]any, error) {
	if filePath == "" && inline == "" {
		return nil, fmt.Errorf("either --file or --data is required")
	}
	var raw []byte
	if inline != "" {
		raw = []byte(inline)
	} else {
		b, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("read file: %w", err)
		}
		raw = b
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}
	return out, nil
}
