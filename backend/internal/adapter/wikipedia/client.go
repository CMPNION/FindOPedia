package wikipedia

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"findopedia/internal/usecase/port"
)

const baseURL = "https://en.wikipedia.org/w/api.php"

type Client struct {
	http *http.Client
}

func NewClient() *Client {
	return &Client{http: &http.Client{}}
}

func (c *Client) FetchRandom(ctx context.Context) (*port.WikiPage, error) {
	url := baseURL + "?action=query&list=random&rnnamespace=0&rnlimit=1&format=json"
	resp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Query struct {
			Random []struct {
				ID    int    `json:"id"`
				Title string `json:"title"`
			} `json:"random"`
		} `json:"query"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode random: %w", err)
	}
	if len(result.Query.Random) == 0 {
		return nil, fmt.Errorf("no random article returned")
	}

	pageID := result.Query.Random[0].ID
	return c.FetchByPageID(ctx, pageID)
}

func (c *Client) FetchByPageID(ctx context.Context, pageID int) (*port.WikiPage, error) {
	url := fmt.Sprintf("%s?action=query&pageids=%d&prop=extracts|info&explaintext=true&exintro=false&format=json", baseURL, pageID)
	resp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Query struct {
			Pages map[string]struct {
				PageID  int    `json:"pageid"`
				Title   string `json:"title"`
				Extract string `json:"extract"`
			} `json:"pages"`
		} `json:"query"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode page: %w", err)
	}

	for _, page := range result.Query.Pages {
		slug := strings.ReplaceAll(page.Title, " ", "_")
		return &port.WikiPage{
			PageID:  page.PageID,
			Title:   page.Title,
			Slug:    slug,
			Extract: page.Extract,
		}, nil
	}
	return nil, fmt.Errorf("page %d not found", pageID)
}

func (c *Client) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "FindOPedia/1.0 (https://github.com/findopedia; contact@findopedia.app)")
	return c.http.Do(req)
}
