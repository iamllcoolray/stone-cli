package updater

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	utiLITIBaseURL         = "https://gurkenlabs.itch.io/litiengine"
	utiLITIFileURLEndpoint = "/file/"
)

type Updater struct {
	http        *http.Client
	installPath string
}

func New(installPath string, httpClient *http.Client) *Updater {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 5 * time.Minute,
		}
	}
	return &Updater{
		installPath: installPath,
		http:        httpClient,
	}
}

func BuildDownloadURL(uploadID string) string {
	return utiLITIBaseURL + utiLITIFileURLEndpoint + uploadID
}

func (u *Updater) Run(downloadURL string, newVersion string) error {
	fmt.Println("Downloading utiLITI...")
	zipPath, err := u.download(downloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	fmt.Println("Extracting...")
	extractedPath, err := u.extract(zipPath)
	if err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	defer u.cleanup(zipPath, extractedPath)

	fmt.Println("Backing up existing installation...")
	if err := u.backup(); err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	fmt.Println("Installing...")
	if err := u.replace(extractedPath); err != nil {
		// replace failed — attempt to restore backup
		backupPath := u.installPath + ".bak"
		if _, statErr := os.Stat(backupPath); statErr == nil {
			fmt.Println("Install failed — restoring backup...")
			if restoreErr := os.Rename(backupPath, u.installPath); restoreErr != nil {
				return fmt.Errorf("install failed and restore failed: %w", err)
			}
			fmt.Println("Backup restored.")
		}
		return fmt.Errorf("install failed: %w", err)
	}

	fmt.Printf("utiLITI %s installed successfully.\n", newVersion)
	return nil
}

func (u *Updater) download(url string) (string, error) {
	resp, err := u.http.Get(url)
	if err != nil {
		return "", fmt.Errorf("downloading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	tmp, err := os.CreateTemp("", "stone-*.zip")
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		return "", fmt.Errorf("writing to temp file: %w", err)
	}

	return tmp.Name(), nil
}

func (u *Updater) extract(zipPath string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("opening zip: %w", err)
	}
	defer r.Close()

	tmpDir, err := os.MkdirTemp("", "stone-extract-*")
	if err != nil {
		return "", fmt.Errorf("creating temp dir: %w", err)
	}

	for _, f := range r.File {
		destPath := filepath.Join(tmpDir, f.Name)

		// guard against zip slip attack
		if !strings.HasPrefix(destPath, filepath.Clean(tmpDir)+string(os.PathSeparator)) {
			return "", fmt.Errorf("invalid file path in zip: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, f.Mode()); err != nil {
				return "", fmt.Errorf("creating directory: %w", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return "", fmt.Errorf("creating parent dir: %w", err)
		}

		if err := extractFile(f, destPath); err != nil {
			return "", err
		}
	}

	return tmpDir, nil
}

func extractFile(f *zip.File, destPath string) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening zip file: %w", err)
	}
	defer rc.Close()

	dest, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, rc); err != nil {
		return fmt.Errorf("extracting file: %w", err)
	}
	return nil
}

func (u *Updater) backup() error {
	if _, err := os.Stat(u.installPath); os.IsNotExist(err) {
		return nil
	}

	backupPath := u.installPath + ".bak"

	if err := os.RemoveAll(backupPath); err != nil {
		return fmt.Errorf("removing old backup: %w", err)
	}

	if err := os.Rename(u.installPath, backupPath); err != nil {
		return fmt.Errorf("creating backup: %w", err)
	}

	return nil
}

func (u *Updater) replace(extractedPath string) error {
	if err := os.MkdirAll(u.installPath, 0755); err != nil {
		return fmt.Errorf("creating install dir: %w", err)
	}

	err := filepath.Walk(extractedPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(extractedPath, path)
		if err != nil {
			return fmt.Errorf("getting relative path: %w", err)
		}

		destPath := filepath.Join(u.installPath, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("creating parent dir: %w", err)
		}

		if err := os.Rename(path, destPath); err != nil {
			if err := copyFile(path, destPath, info.Mode()); err != nil {
				return fmt.Errorf("copying file: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("replacing files: %w", err)
	}
	return nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}

func (u *Updater) cleanup(tmpPaths ...string) error {
	for _, path := range tmpPaths {
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("cleaning up %s: %w", path, err)
		}
	}

	backupPath := u.installPath + ".bak"
	if _, err := os.Stat(backupPath); err == nil {
		if err := os.RemoveAll(backupPath); err != nil {
			return fmt.Errorf("removing backup: %w", err)
		}
	}

	return nil
}
