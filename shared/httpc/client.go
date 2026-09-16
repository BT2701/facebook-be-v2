package httpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout = 3 * time.Second
	maxAttempts    = 3
)

var client = &http.Client{Timeout: defaultTimeout}

func GetJSON(ctx context.Context, rawURL string, dest interface{}) error {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if attempt > 1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(attempt-1) * 150 * time.Millisecond):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		_ = resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			continue
		}

		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("upstream %s: %s", rawURL, resp.Status)
			continue
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf("upstream %s: %s", rawURL, resp.Status)
		}
		if dest == nil {
			return nil
		}
		if err := json.Unmarshal(body, dest); err != nil {
			return err
		}
		return nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("upstream %s: request failed", rawURL)
	}
	return lastErr
}
