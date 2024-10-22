package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

// curl -X GET 'https://registry.hub.docker.com/v2/repositories/{namespace}/{repository}/tags'

func main() {
	// TODO
	fmt.Println("Hello, World!")

	fp, err := os.Open("test.json")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	httpClient := &http.Client{}
	api := NewDockerHubAPI("https://registry.hub.docker.com", httpClient)

	response, err := api.ListRepositoryTags("awsguru", "aws-lambda-adapter")
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	logger.Info("response", zap.Any("response", response))
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type DockerHubAPI struct {
	HTTPClient HTTPClient
	BaseURL    string
}

func NewDockerHubAPI(baseURL string, httpClient HTTPClient) *DockerHubAPI {
	return &DockerHubAPI{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
	}
}

func (api *DockerHubAPI) ListRepositoryTags(namespace, repository string) (*ListTagsResponse, error) {
	path := fmt.Sprintf("/v2/repositories/%s/%s/tags", namespace, repository)
	request, err := http.NewRequest(http.MethodGet, api.BaseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	defer request.Body.Close()
	var response ListTagsResponse

	if err := json.NewDecoder(request.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

type ListTagsResponse struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous *string  `json:"previous"`
	Results  []Result `json:"results"`
}

type Result struct {
	Creator             int       `json:"creator"`
	ID                  int       `json:"id"`
	Images              []Image   `json:"images"`
	LastUpdated         time.Time `json:"last_updated"`
	LastUpdater         int       `json:"last_updater"`
	LastUpdaterUsername string    `json:"last_updater_username"`
	Name                string    `json:"name"`
	Repository          int       `json:"repository"`
	FullSize            int       `json:"full_size"`
	V2                  bool      `json:"v2"`
	TagStatus           string    `json:"tag_status"`
	TagLastPulled       time.Time `json:"tag_last_pulled"`
	TagLastPushed       time.Time `json:"tag_last_pushed"`
	MediaType           string    `json:"media_type"`
	ContentType         string    `json:"content_type"`
	Digest              string    `json:"digest"`
}

type Image struct {
	Architecture string    `json:"architecture"`
	Features     string    `json:"features"`
	Variant      *string   `json:"variant"`
	Digest       string    `json:"digest"`
	OS           string    `json:"os"`
	OSFeatures   string    `json:"os_features"`
	OSVersion    *string   `json:"os_version"`
	Size         int       `json:"size"`
	Status       string    `json:"status"`
	LastPulled   time.Time `json:"last_pulled"`
	LastPushed   time.Time `json:"last_pushed"`
}
