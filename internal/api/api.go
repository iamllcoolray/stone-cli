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

// Upload represents a single itch.io upload entry
type Upload struct {
	ID          int    `json:"id"`
	Filename    string `json:"filename"`
	PWindows    bool   `json:"p_windows"`
	PLinux      bool   `json:"p_linux"`
	POSX        bool   `json:"p_osx"`
	DisplayName string `json:"display_name"`
	BuildID     int    `json:"build_id"`
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

// FetchLatestVersion returns the version string for the current platform upload
func (c *Client) FetchLatestVersion() (string, error) {
	upload, err := c.FetchPlatformUpload()
	if err != nil {
		return "", err
	}

	if upload.BuildID == 0 {
		return "", fmt.Errorf("no build found for upload %d", upload.ID)
	}

	url := fmt.Sprintf("%s/%s/build/%d/info", itchAPIBaseURL, c.apiKey, upload.BuildID)

	resp, err := c.http.Get(url)
	if err != nil {
		return "", fmt.Errorf("fetching build info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result buildResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding build info: %w", err)
	}

	if result.Build.UserVersion == "" {
		return "", fmt.Errorf("version not found in build info")
	}

	return result.Build.UserVersion, nil
}

// FetchDownloadURL returns a direct download URL for the given upload ID
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
