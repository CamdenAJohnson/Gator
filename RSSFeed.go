package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var Feed RSSFeed

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create http request: %v\n", err)
	}

	request.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to complete http request: %v\n", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v\n", err)
	}

	if err := xml.Unmarshal(body, &Feed); err != nil {
		return nil, fmt.Errorf("Failed to parse response: %v\n", err)
	}

	Feed.Channel.Title = html.UnescapeString(Feed.Channel.Title)
	Feed.Channel.Description = html.UnescapeString(Feed.Channel.Description)

	return &Feed, nil
}