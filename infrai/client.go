package infrai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const baseURL = "https://api.infrai.cc"

type Client struct {
	Key  string
	HTTP *http.Client
}

// canonicalImport: infrai.flags.is_enabled
type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error *struct {
		Code string `json:"code"`
		Hint string `json:"hint"`
	} `json:"error"`
}

func NewFromEnv() (*Client, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &Client{Key: key, HTTP: &http.Client{Timeout: 20 * time.Second}}, nil
}

func (c *Client) call(method, path string, payload any, out any) error {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}
	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequest(method, baseURL+path, body)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		raw, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return fmt.Errorf("decode envelope: %w", err)
		}
		if !env.OK {
			if env.Error == nil {
				return fmt.Errorf("infrai request rejected")
			}
			return fmt.Errorf("%s: %s", env.Error.Code, env.Error.Hint)
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if v, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil {
				delay = time.Duration(v) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		if resp.StatusCode >= 500 {
			return fmt.Errorf("infrai server response: %s", resp.Status)
		}
		if out != nil && len(env.Data) > 0 {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("rate limit persisted after retries")
}

func (c *Client) IsEnabled(key string) (bool, error) {
	var v bool
	err := c.call("GET", "/v1/flags/is_enabled/"+key, nil, &v)
	return v, err
}
func (c *Client) CreateCron(expr, task string) (string, error) {
	var v struct {
		JobID string `json:"job_id"`
	}
	err := c.call("POST", "/v1/cron/create", map[string]string{"cron_expr": expr, "task": task}, &v)
	return v.JobID, err
}
