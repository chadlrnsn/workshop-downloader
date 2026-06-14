package steamcmd

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"workshop-downloader/internal/domain"
)

const steamCmdZipURL = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip"

// VerifyOrInstall checks if steamcmd.exe exists at steamCmdPath; if not, downloads and extracts it.
func VerifyOrInstall(ctx context.Context, steamCmdPath string, logFn OutputHandler) error {
	dir := filepath.Dir(steamCmdPath)
	
	// 1. Check if steamcmd.exe already exists
	if _, err := os.Stat(steamCmdPath); err == nil {
		logFn(domain.LogEvent{
			Message:   "SteamCMD executable found: " + steamCmdPath,
			Timestamp: time.Now(),
		})
		return nil
	}

	logFn(domain.LogEvent{
		Message:   "SteamCMD not found. Starting automatic installation in: " + dir,
		Timestamp: time.Now(),
	})

	// Create directory structure
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	zipPath := filepath.Join(dir, "steamcmd.zip")

	// 2. Download SteamCMD installer zip
	logFn(domain.LogEvent{
		Message:   "Downloading SteamCMD from " + steamCmdZipURL + "...",
		Timestamp: time.Now(),
	})
	
	if err := downloadFile(ctx, steamCmdZipURL, zipPath); err != nil {
		return fmt.Errorf("failed to download SteamCMD: %w", err)
	}

	logFn(domain.LogEvent{
		Message:   "Download complete. Extracting files...",
		Timestamp: time.Now(),
	})

	// 3. Extract Zip file
	if err := unzip(zipPath, dir); err != nil {
		_ = os.Remove(zipPath)
		return fmt.Errorf("failed to extract zip archive: %w", err)
	}

	// Clean up zip file
	_ = os.Remove(zipPath)

	logFn(domain.LogEvent{
		Message:   "Extraction successful. Running initial SteamCMD update sequence (this may take a minute)...",
		Timestamp: time.Now(),
	})

	// 4. Run steamcmd once with +quit to force built-in update
	cmd := exec.CommandContext(ctx, steamCmdPath, "+quit")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x00000008,
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start initial steamcmd self-update: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	logPipe := func(r io.Reader, isErr bool) {
		defer wg.Done()
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			logFn(domain.LogEvent{
				Message:   scanner.Text(),
				Timestamp: time.Now(),
				IsError:   isErr,
			})
		}
	}

	go logPipe(stdout, false)
	go logPipe(stderr, true)

	wg.Wait()
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to run initial steamcmd self-update: %w", err)
	}

	logFn(domain.LogEvent{
		Message:   "SteamCMD is fully updated and ready for operations.",
		Timestamp: time.Now(),
	})

	return nil
}

func downloadFile(ctx context.Context, url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
