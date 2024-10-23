package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

// curl -X GET 'https://registry.hub.docker.com/v2/repositories/{namespace}/{repository}/tags'

func main() {
	args := os.Args
	namespace, repository, err := ParseArgs(args)
	if err != nil {
		fmt.Println(Usage())
		os.Exit(1)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	httpClient := &http.Client{}
	api := NewDockerHubAPI(logger, "https://registry.hub.docker.com", httpClient)

	response, err := api.ListRepositoryTags(namespace, repository)
	if err != nil {
		if errors.Is(err, RepositoryNotFound) {
			fmt.Println("repository not found")
			os.Exit(1)
		}
		panic(err)
	}

	logger.Info("response", zap.Any("response", response))
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type DockerHubAPI struct {
	logger     *zap.Logger
	HTTPClient HTTPClient
	BaseURL    string
}

func NewDockerHubAPI(logger *zap.Logger, baseURL string, httpClient HTTPClient) *DockerHubAPI {
	return &DockerHubAPI{
		logger:     logger,
		BaseURL:    baseURL,
		HTTPClient: httpClient,
	}
}

func Usage() string {
	return "Usage: $0 <namespace>/<repository>"
}

func ParseArgs(args []string) (string, string, error) {
	if len(args) != 2 {
		return "", "", fmt.Errorf("invalid arguments")
	}
	a := strings.Split(args[1], "/")
	if len(a) != 2 {
		return "", "", fmt.Errorf("invalid arguments")
	}
	return a[0], a[1], nil
}

var RepositoryNotFound = fmt.Errorf("repository not found")

func (api *DockerHubAPI) ListRepositoryTags(namespace, repository string) (*ListTagsResponse, error) {
	url := fmt.Sprintf("%s/v2/repositories/%s/%s/tags", api.BaseURL, namespace, repository)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := api.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, RepositoryNotFound
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var response ListTagsResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
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
