// Package goindex is an index.golang.org compatible client.
package goindex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Client is the main entry point to an index
// service. The zero value is ready to use.
type Client struct {
	http *http.Client
	url  string
}

// Get gets a list of module versions from the index service.
// Since is the ?since=<time> parameter and limit is the ?limit=<int>
// parameter. See https://index.golang.org for more info.
func (c *Client) Get(ctx context.Context, since time.Time, limit int) (ModuleVersions, error) {
	h := c.http
	if h == nil {
		h = http.DefaultClient
	}
	u := c.url
	if u == "" {
		u = "https://index.golang.org/index"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	q := req.URL.Query()
	if !since.IsZero() {
		q.Set("since", since.Format(time.RFC3339))
	}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	req.URL.RawQuery = q.Encode()
	resp, err := h.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var mods ModuleVersions
	dec := json.NewDecoder(resp.Body)
	for dec.More() {
		var m ModuleVersion
		err = dec.Decode(&m)
		if err != nil {
			return nil, fmt.Errorf("json.Decode: %w", err)
		}
		mods = append(mods, &m)
	}
	return mods, nil
}

// ModuleVersions is a paginated list
// of module versions that knows how to
// get the next version
type ModuleVersions []*ModuleVersion

// Next gets the next times right after the last module version
// in the slice.
func (ms ModuleVersions) Next(ctx context.Context, c *Client, limit int) (ModuleVersions, error) {
	if len(ms) == 0 {
		return nil, io.EOF
	}
	since := ms[len(ms)-1].Timestamp
	mods, err := c.Get(ctx, since, limit)
	if err != nil {
		return nil, err
	}
	if len(mods) <= 1 {
		return nil, io.EOF
	}
	return mods[1:], nil
}

// ModuleVersion represents a single entry
// line in an index api response
type ModuleVersion struct {
	Path      string
	Version   string
	Timestamp time.Time
}
