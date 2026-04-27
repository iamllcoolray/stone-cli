package scraper

import "net/http"

const (
	utiLITIPageURL             = "https://gurkenlabs.itch.io/litiengine"
	utiLITIDownloadURLEndpoint = "https://gurkenlabs.itch.io/litiengine/download"
)

type Client struct {
	http *http.Client
}

func New() *Client {
	return &Client{}
}

func (c *Client) FetchLatestVersion() (string, error) {
	return "", nil
}

func (c *Client) fetchCSRFToken() (string, error) {
	return "", nil
}

func (c *Client) FetchDownloadURL() (string, error) {
	return "", nil
}

func (c *Client) FetchPlatformUploadID(downloadPageURL string) (string, error) {
	return "", nil
}

func platformIcon() string {
	return ""
}
