package dune

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

var ErrorReqUnsuccessful = errors.New("request was not successful")

type ErrorResponse struct {
	Error string `json:"error"`
}

type RateLimit struct {
	Limit     int
	Remaining int
	Reset     int64
}

type APIError struct {
	StatusCode  int
	StatusText  string
	BodySnippet string
	RateLimit   *RateLimit
	RetryAfter  time.Duration
}

func (e *APIError) Error() string {
	if e.BodySnippet != "" {
		return fmt.Sprintf("http %d %s: %s", e.StatusCode, e.StatusText, e.BodySnippet)
	}
	return fmt.Sprintf("http %d %s", e.StatusCode, e.StatusText)
}

type RetryPolicy struct {
	MaxAttempts           int
	InitialBackoff        time.Duration
	MaxBackoff            time.Duration
	Jitter                time.Duration
	RetryableStatusCodes  []int
}

var defaultRetryPolicy = RetryPolicy{
	MaxAttempts:          3,
	InitialBackoff:       500 * time.Millisecond,
	MaxBackoff:           5 * time.Second,
	Jitter:               100 * time.Millisecond,
	RetryableStatusCodes: []int{429, 500, 502, 503, 504},
}

func parseRateLimitHeaders(h http.Header) *RateLimit {
	limStr := h.Get("X-RateLimit-Limit")
	remStr := h.Get("X-RateLimit-Remaining")
	resetStr := h.Get("X-RateLimit-Reset")

	var lim, rem int
	var reset int64

	if limStr != "" {
		if v, err := strconv.Atoi(limStr); err == nil {
			lim = v
		}
	}
	if remStr != "" {
		if v, err := strconv.Atoi(remStr); err == nil {
			rem = v
		}
	}
	if resetStr != "" {
		if v, err := strconv.ParseInt(resetStr, 10, 64); err == nil {
			reset = v
		}
	}

	if lim == 0 && rem == 0 && reset == 0 {
		return nil
	}
	return &RateLimit{Limit: lim, Remaining: rem, Reset: reset}
}

func nextBackoff(attempt int, p RetryPolicy) time.Duration {
	b := p.InitialBackoff
	for i := 1; i < attempt; i++ {
		b *= 2
		if b > p.MaxBackoff {
			b = p.MaxBackoff
			break
		}
	}
	if p.Jitter > 0 {
		b += p.Jitter
	}
	return b
}

func decodeBody(resp *http.Response, dest interface{}) error {
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(dest)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	return nil
}

func httpRequest(apiKey string, req *http.Request) (*http.Response, error) {
	req.Header.Add("X-DUNE-API-KEY", apiKey)
	p := defaultRetryPolicy
	attempt := 1
	for {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if attempt >= p.MaxAttempts {
				return nil, fmt.Errorf("failed to send request: %w", err)
			}
			time.Sleep(nextBackoff(attempt, p))
			attempt++
			continue
		}

		if resp.StatusCode == 200 {
			return resp, nil
		}

		defer resp.Body.Close()
		snippetBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		var er ErrorResponse
		_ = json.Unmarshal(snippetBytes, &er)
		msg := string(snippetBytes)
		if er.Error != "" {
			msg = er.Error
		}
		rl := parseRateLimitHeaders(resp.Header)
		retryAfter := time.Duration(0)
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, err := strconv.Atoi(ra); err == nil {
				retryAfter = time.Duration(secs) * time.Second
			}
		}
		apiErr := &APIError{
			StatusCode:  resp.StatusCode,
			StatusText: resp.Status,
			BodySnippet: msg,
			RateLimit:  rl,
			RetryAfter: retryAfter,
		}
		retryable := false
		for _, code := range p.RetryableStatusCodes {
			if resp.StatusCode == code {
				retryable = true
				break
			}
		}
		if retryable && attempt < p.MaxAttempts {
			sleep := nextBackoff(attempt, p)
			if apiErr.RetryAfter > 0 && apiErr.RetryAfter > sleep {
				sleep = apiErr.RetryAfter
			}
			time.Sleep(sleep)
			attempt++
			continue
		}
		return nil, fmt.Errorf("%w: %v", ErrorReqUnsuccessful, apiErr)
	}
}
