package steamcmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"workshop-downloader/internal/domain"
)

type OutputHandler func(event domain.LogEvent)
type ProgressHandler func(progress float64)
type ServiceCodePrompt struct {
	PromptType string // "email" or "2fa"
	Response   chan string
}

type SteamCmdRunner struct {
	mu           sync.Mutex
	steamCmdPath string
	codePrompt   chan*ServiceCodePrompt
}

func NewSteamCmdRunner(steamCmdPath string) *SteamCmdRunner {
	return &SteamCmdRunner{
		steamCmdPath: steamCmdPath,
		codePrompt:   make(chan *ServiceCodePrompt),
	}
}

func (r *SteamCmdRunner) UpdatePath(newPath string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.steamCmdPath = newPath
}

// CodePromptChan exposes the channel for Wails backend to catch auth code requests
func (r *SteamCmdRunner) CodePromptChan() <-chan *ServiceCodePrompt {
	return r.codePrompt
}

// RunSteamCmd executes steamcmd with custom arguments and handles prompts
func (r *SteamCmdRunner) RunSteamCmd(
	ctx context.Context,
	args []string,
	logHandler OutputHandler,
	progressHandler ProgressHandler,
) error {
	r.mu.Lock()
	cmdPath := r.steamCmdPath
	r.mu.Unlock()

	cmd := exec.CommandContext(ctx, cmdPath, args...)
	
	// Crucial for Windows: hide console window
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x00000008, // DETACHED_PROCESS or CREATE_NO_WINDOW
	}

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	defer stdinPipe.Close()

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start steamcmd: %w", err)
	}

	// Regexp patterns
	progressRegex := regexp.MustCompile(`Update state \((.*?)\) downloading, progress: (\d+\.\d+)`)
	steamGuardEmailRegex := regexp.MustCompile(`(?i)(Steam Guard code|Enter the Steam Guard code that was sent to your email address)`)
	steamGuard2FARegex := regexp.MustCompile(`(?i)(Two-factor code|Enter the Steam Guard code from your authenticator)`)
	errorRegex := regexp.MustCompile(`(?i)(ERROR!|Failed to|Failed auth|Invalid Password|Login Failed)`)

	var runnerErr error
	var wg sync.WaitGroup
	wg.Add(2)

	// Stream logs/stderr in background
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			logHandler(domain.LogEvent{
				Message:   line,
				Timestamp: time.Now(),
				IsError:   true,
			})
		}
	}()

	// Scan stdout and state machine
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdoutPipe)
		var lineBuffer []byte

		for {
			b, err := reader.ReadByte()
			if err != nil {
				if err != io.EOF && len(lineBuffer) > 0 {
					logLine(string(lineBuffer), logHandler, progressHandler, progressRegex, steamGuardEmailRegex, steamGuard2FARegex, errorRegex, stdinPipe, r.codePrompt)
				}
				break
			}

			if b == '\n' || b == '\r' {
				if len(lineBuffer) > 0 {
					logLine(string(lineBuffer), logHandler, progressHandler, progressRegex, steamGuardEmailRegex, steamGuard2FARegex, errorRegex, stdinPipe, r.codePrompt)
					lineBuffer = lineBuffer[:0]
				}
				continue
			}

			lineBuffer = append(lineBuffer, b)

			// Check for interactive prompts which do not emit a trailing newline
			currentStr := string(lineBuffer)
			if strings.Contains(currentStr, "Steam Guard code:") || strings.Contains(currentStr, "Two-factor code:") {
				logLine(currentStr, logHandler, progressHandler, progressRegex, steamGuardEmailRegex, steamGuard2FARegex, errorRegex, stdinPipe, r.codePrompt)
				lineBuffer = lineBuffer[:0]
			}
		}
	}()

	// Wait for process to exit or context cancellation
	errChan := make(chan error, 1)
	go func() {
		wg.Wait()
		errChan <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// Gracefully terminate, then force kill
		_ = cmd.Process.Signal(os.Interrupt)
		time.AfterFunc(2*time.Second, func() {
			_ = cmd.Process.Kill()
		})
		return ctx.Err()
	case err := <-errChan:
		if err != nil {
			if runnerErr != nil {
				return runnerErr
			}
			return fmt.Errorf("steamcmd process returned error: %w", err)
		}
	}

	return nil
}

// DownloadItem executes steamcmd to download a workshop item
func (r *SteamCmdRunner) DownloadItem(
	ctx context.Context,
	appID string,
	workshopID string,
	username string,
	password string,
	logHandler OutputHandler,
	progressHandler ProgressHandler,
) error {
	args := []string{}
	if username != "" && username != "anonymous" {
		args = append(args, "+login", username)
		if password != "" {
			args = append(args, password)
		}
	} else {
		args = append(args, "+login", "anonymous")
	}

	args = append(args, "+workshop_download_item", appID, workshopID, "validate", "+quit")
	return r.RunSteamCmd(ctx, args, logHandler, progressHandler)
}

// Login performs a blocking login to cache credentials / Sentry ticket
func (r *SteamCmdRunner) Login(
	ctx context.Context,
	username string,
	password string,
	logHandler OutputHandler,
) error {
	args := []string{}
	if username != "" && username != "anonymous" {
		args = append(args, "+login", username)
		if password != "" {
			args = append(args, password)
		}
	} else {
		args = append(args, "+login", "anonymous")
	}
	args = append(args, "+quit")

	// Pass no-op progress handler for login
	dummyProgress := func(prog float64) {}
	return r.RunSteamCmd(ctx, args, logHandler, dummyProgress)
}

func logLine(
	line string,
	logHandler OutputHandler,
	progressHandler ProgressHandler,
	progressRegex *regexp.Regexp,
	steamGuardEmailRegex *regexp.Regexp,
	steamGuard2FARegex *regexp.Regexp,
	errorRegex *regexp.Regexp,
	stdin io.Writer,
	codePromptChan chan *ServiceCodePrompt,
) {
	cleanLine := strings.TrimSpace(line)
	if cleanLine == "" {
		return
	}

	// 1. Emit live logs to UI
	logHandler(domain.LogEvent{
		Message:   cleanLine,
		Timestamp: time.Now(),
		IsError:   errorRegex.MatchString(cleanLine),
	})

	// 2. Parse Progress
	if progressRegex != nil {
		if matches := progressRegex.FindStringSubmatch(cleanLine); len(matches) > 1 {
			if progressVal, err := strconv.ParseFloat(matches[2], 64); err == nil {
				progressHandler(progressVal)
			}
		}
	}

	// 3. Steam Guard Interaction
	var promptType string
	if steamGuardEmailRegex.MatchString(cleanLine) {
		promptType = "email"
	} else if steamGuard2FARegex.MatchString(cleanLine) {
		promptType = "2fa"
	}

	if promptType != "" {
		prompt := &ServiceCodePrompt{
			PromptType: promptType,
			Response:   make(chan string),
		}
		// Send prompt to orchestrator (app bindings)
		select {
		case codePromptChan <- prompt:
			// Wait for reply from the GUI binding
			code := <-prompt.Response
			if code != "" {
				_, _ = io.WriteString(stdin, code+"\n")
			} else {
				// Abort
				_ = stdin.(*os.File).Close()
			}
		case <-time.After(3 * time.Minute):
			// Timeout, do nothing
			_ = stdin.(*os.File).Close()
		}
	}
}
