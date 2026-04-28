package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	utiLITIBaseURL             = "https://gurkenlabs.itch.io/litiengine"
	utiLITIDownloadURLEndpoint = "/download_url"
)

type userAgentTransport struct {
	agent string
	base  http.RoundTripper
}

func (t *userAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", t.agent)
	return t.base.RoundTrip(r)
}

type Client struct {
	http *http.Client
}

func New() *Client {
	return &Client{
		http: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &userAgentTransport{
				agent: "stone-cli/1.0",
				base:  http.DefaultTransport,
			},
		},
	}
}

func (c *Client) FetchLatestVersion() (string, error) {
	resp, err := c.http.Get(utiLITIBaseURL)
	if err != nil {
		return "", fmt.Errorf("fetching page: %w", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("parsing html: %w", err)
	}

	version := findLatestVersion(doc)
	if version == "" {
		return "", fmt.Errorf("latest version not found")
	}
	return version, nil
}

func findLatestVersion(n *html.Node) string {
	// find the game_devlog section
	if n.Type == html.ElementNode && n.Data == "section" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == "game_devlog" {
				// grab the first <a> inside it
				link := findFirstLink(n)
				if link != "" {
					return strings.TrimPrefix(link, "Release ")
				}
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if v := findLatestVersion(child); v != "" {
			return v
		}
	}
	return ""
}

func findFirstLink(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "a" {
		// get the text content of the link
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			return n.FirstChild.Data
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if text := findFirstLink(child); text != "" {
			return text
		}
	}
	return ""
}
func (c *Client) fetchCSRFToken() (string, error) {
	resp, err := c.http.Get(utiLITIBaseURL)
	if err != nil {
		return "", fmt.Errorf("fetching page: %w", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("parsing html: %w", err)
	}

	token := findCSRFToken(doc)
	if token == "" {
		return "", fmt.Errorf("csrf token not found")
	}
	return token, nil
}

func findCSRFToken(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "meta" {
		var name, value string
		for _, attr := range n.Attr {
			switch attr.Key {
			case "name":
				name = attr.Val
			case "value":
				value = attr.Val
			}
		}
		if name == "csrf_token" {
			return value
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if token := findCSRFToken(child); token != "" {
			return token
		}
	}
	return ""
}

type downloadURLResponse struct {
	URL string `json:"url"`
}

func (c *Client) FetchDownloadURL() (string, error) {
	token, err := c.fetchCSRFToken()
	if err != nil {
		return "", fmt.Errorf("fetching csrf token: %w", err)
	}

	req, err := http.NewRequest("POST", utiLITIBaseURL+utiLITIDownloadURLEndpoint, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("X-CSRF-Token", token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("posting download url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result downloadURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	if result.URL == "" {
		return "", fmt.Errorf("download url not found in response")
	}
	return result.URL, nil
}

func (c *Client) FetchPlatformUploadID(downloadPageURL string) (string, error) {
	resp, err := c.http.Get(downloadPageURL)
	if err != nil {
		return "", fmt.Errorf("fetching download page: %w", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("parsing html: %w", err)
	}

	icon := platformIcon()
	if icon == "" {
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	uploadID := findUploadID(doc, icon)
	if uploadID == "" {
		return "", fmt.Errorf("upload id not found for platform: %s", runtime.GOOS)
	}
	return uploadID, nil
}

func findUploadID(n *html.Node, icon string) string {
	// find every <div class="upload">
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == "upload" {
				// check if this upload div contains the platform icon
				if hasIcon(n, icon) {
					// grab data-upload_id from the download_btn
					return findDataUploadID(n)
				}
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if id := findUploadID(child, icon); id != "" {
			return id
		}
	}
	return ""
}

func hasIcon(n *html.Node, icon string) bool {
	// walk the upload div looking for a <span> whose class contains the icon string
	if n.Type == html.ElementNode && n.Data == "span" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, icon) {
				return true
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if hasIcon(child, icon) {
			return true
		}
	}
	return false
}

func findDataUploadID(n *html.Node) string {
	// find the <a class="button download_btn"> and return its data-upload_id
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, "download_btn") {
				for _, attr := range n.Attr {
					if attr.Key == "data-upload_id" {
						return attr.Val
					}
				}
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if id := findDataUploadID(child); id != "" {
			return id
		}
	}
	return ""
}

func platformIcon() string {
	switch runtime.GOOS {
	case "linux":
		return "icon-tux"
	case "windows":
		return "icon-windows8"
	case "darwin":
		return "icon-apple"
	default:
		return ""
	}
}
