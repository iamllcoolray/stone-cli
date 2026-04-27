package updater

import "net/http"

type Updater struct {
	http        *http.Client
	installPath string
}

func New(installPath string) *Updater {
	return &Updater{}
}

func (u *Updater) Run(uploadID string, newVersion string) error {
	return nil
}

func buildDownloadURL(uploadID string) string {
	return ""
}

func (u *Updater) download(url string) (string, error) {
	return "", nil
}

func (u *Updater) extract(zipPath string) (string, error) {
	return "", nil
}

func (u *Updater) replace(extractedPath string) error {
	return nil
}

func (u *Updater) backup() error {
	return nil
}

func (u *Updater) cleanup(tmpPaths ...string) error {
	return nil
}
