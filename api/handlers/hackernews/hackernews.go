package hackernews

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"peek/api/httpclient"
	"time"
)

type HackerNewsClient struct {
	client *httpclient.Client
}

func NewHackerNewsClient() *HackerNewsClient {
	baseURL := "https://hacker-news.firebaseio.com/v0/"
	client := httpclient.NewClient(baseURL, httpclient.WithTimeout(15*time.Second))
	return &HackerNewsClient{client: client}
}

func (hn *HackerNewsClient) GetTopStories(ctx context.Context) ([]int, error) {
	var ids []int
	err := hn.client.Get(ctx, "topstories.json?print=pretty", &ids)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch top stories: %w", err)
	}
	return ids, nil
}

func (hn *HackerNewsClient) GetStory(ctx context.Context, id int) (map[string]interface{}, error) {
	var story map[string]interface{}
	endpoint := fmt.Sprintf("item/%d.json?print=pretty", id)
	err := hn.client.Get(ctx, endpoint, &story)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch story with ID %d: %w", id, err)
	}
	return story, nil
}

func HackerNewsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	hn := NewHackerNewsClient()

	topStoryIDs, err := hn.GetTopStories(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching top stories: %v", err), http.StatusInternalServerError)
		return
	}

	stories := make([]map[string]interface{}, 0, 10)
	for _, id := range topStoryIDs[:10] {
		story, err := hn.GetStory(ctx, id)
		if err != nil {
			continue
		}
		stories = append(stories, story)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stories)
}
