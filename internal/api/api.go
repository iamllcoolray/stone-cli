package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

const (
	utiLITIGameID  = "731589"
	itchAPIBaseURL = "https://itch.io/api/1"
)

type Client struct {
	http   *http.Client
	apiKey string
}

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) HTTPClient() *http.Client {
	return c.http
}

type Build struct {
	ID          int    `json:"id"`
	UserVersion string `json:"user_version"`
}

type Upload struct {
	ID          int    `json:"id"`
	Filename    string `json:"filename"`
	DisplayName string `json:"display_name"`
	ChannelName string `json:"channel_name"`
	PWindows    bool   `json:"p_windows"`
	PLinux      bool   `json:"p_linux"`
	POSX        bool   `json:"p_osx"`
	Build       Build  `json:"build"`
}

type uploadsResponse struct {
	Uploads []Upload `json:"uploads"`
}

type downloadResponse struct {
	URL string `json:"url"`
}

type buildResponse struct {
	Build struct {
		UserVersion string `json:"user_version"`
	} `json:"build"`
}

// FetchPlatformUpload returns the upload matching the current platform
func (c *Client) FetchPlatformUpload() (*Upload, error) {
	url := fmt.Sprintf("%s/%s/game/%s/uploads", itchAPIBaseURL, c.apiKey, utiLITIGameID)

	resp, err := c.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching uploads: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result uploadsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding uploads: %w", err)
	}

	for _, u := range result.Uploads {
		if matchesPlatform(u) {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("no upload found for platform: %s", runtime.GOOS)
}

func (c *Client) FetchLatestVersion() (string, error) {
	upload, err := c.FetchPlatformUpload()
	if err != nil {
		return "", err
	}
	if upload.Build.UserVersion == "" {
		return "", fmt.Errorf("version not found for upload %d", upload.ID)
	}
	return upload.Build.UserVersion, nil
}

func (c *Client) FetchDownloadURL(uploadID int) (string, error) {
	url := fmt.Sprintf("%s/%s/upload/%d/download", itchAPIBaseURL, c.apiKey, uploadID)

	resp, err := c.http.Get(url)
	if err != nil {
		return "", fmt.Errorf("fetching download url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result downloadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding download url: %w", err)
	}

	if result.URL == "" {
		return "", fmt.Errorf("download url not found")
	}

	return result.URL, nil
}

// matchesPlatform returns true if the upload matches the current OS
func matchesPlatform(u Upload) bool {
	switch runtime.GOOS {
	case "linux":
		return u.PLinux
	case "windows":
		return u.PWindows
	case "darwin":
		return u.POSX
	default:
		return false
	}
}
