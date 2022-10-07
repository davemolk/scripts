package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// getURLs makes a request to the Wayback Machine and unmarshals
// the resulting response body into a slice of string slices. This
// is returned, along with any error.
func (t *tas) getURLs(url string, timeout int) ([][]string, error) {
	u := t.makeURL(url)
	resp, err := t.makeRequest(u, timeout)
	if err != nil {
		return nil, fmt.Errorf("makeRequest err: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body read err: %w", err)
	}

	var s [][]string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, fmt.Errorf("unmarshal err: %w", err)
	}

	return s, nil
}

// testURLs takes in a URL, makes a request, and
// writes the results to the statusMap.
func (t *tas) testURLs(url string) {
	resp, err := t.makeRequest(url, t.config.timeout)
	if err != nil {
		log.Printf("error in testURLs: %v", err)
		return
	}
	defer resp.Body.Close()
	t.results.add(resp.StatusCode, url)
}

// makeClient returns a single client for reuse
func (t *tas) makeClient(redirect bool) *http.Client {
	return &http.Client{
		CheckRedirect: t.allowRedirects(redirect),
	}
}

// allowRedirects takes in a boolean whose value is determined by the redirects flag
// If the flag is true, allowRedirects returns nil and redirects will be allowed. Otherwise,
// allowRedirects returns a function for the CheckRedirect field of the http.Client that
// blocks redirects.
func (t *tas) allowRedirects(redirect bool) func(*http.Request, []*http.Request) error {
	if redirect {
		return nil
	} else {
		return func(req *http.Request, via []*http.Request) error {
			log.Printf("blocked attempted redirect to %s\n", req.URL.String())
			return http.ErrUseLastResponse
		}
	}
}

// makeURL takes in the user-supplied URL and builds the
// query URL to send to the Wayback Machine.
func (t *tas) makeURL(url string) string {
	now := time.Now()
	curr := now.UnixMilli()
	const begin = "https://web.archive.org/web/timemap/json?url="
	const mid = "&matchType=prefix&collapse=urlkey&output=json&fl=original%2Cmimetype%2Ctimestamp%2Cendtimestamp%2Cgroupcount%2Cuniqcount&filter=!statuscode%3A%5B45%5D..&limit=10000&_="
	u := fmt.Sprintf("%s%s%s%d", begin, url, mid, curr)
	return u
}

// makeRequest takes in a URL and a timeout and returns a response and any error.
func (t *tas) makeRequest(url string, timeout int) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	uA := t.randomUA()
	req.Header.Set("User-Agent", uA)

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// randomUA assembles a slice of 10 user-agents and picks one.
func (t *tas) randomUA() string {
	userAgents := t.getUA()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	agent := r.Intn(len(userAgents))
	return userAgents[agent]
}

// getUA returns a slice of 10 user-agents.
func (t *tas) getUA() []string {
	return []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4692.56 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4889.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/601.7.7 (KHTML, like Gecko) Version/9.1.2 Safari/601.7.7",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36",
	}
}
