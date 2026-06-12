package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"

	"workshop-downloader/internal/config"
	"workshop-downloader/internal/domain"
	"workshop-downloader/internal/downloader"
	"workshop-downloader/internal/steamcmd"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx           context.Context
	cfgManager    *config.ConfigManager
	cmdRunner     *steamcmd.SteamCmdRunner
	downManager   *downloader.DownloadManager
	mu            sync.Mutex
	activePrompt  *steamcmd.ServiceCodePrompt
}

type WailsBroadcaster struct {
	ctx context.Context
}

func (wb *WailsBroadcaster) Emit(eventName string, data interface{}) {
	if wb.ctx != nil {
		runtime.EventsEmit(wb.ctx, eventName, data)
	}
}

func NewApp() *App {
	cfgManager := config.NewConfigManager()
	cfg := cfgManager.GetConfig()
	cmdRunner := steamcmd.NewSteamCmdRunner(cfg.SteamCmdPath)

	return &App{
		cfgManager: cfgManager,
		cmdRunner:  cmdRunner,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	broadcaster := &WailsBroadcaster{ctx: ctx}
	a.downManager = downloader.NewDownloadManager(a.cmdRunner, a.cfgManager, broadcaster)
	a.downManager.Start()

	// Listen for Steam Guard prompts
	go a.listenForAuthPrompts()

	// Check SteamCMD installation on startup
	go a.checkSteamCmdOnStartup()
}

func (a *App) shutdown(ctx context.Context) {
	if a.downManager != nil {
		a.downManager.Stop()
	}
}

func (a *App) listenForAuthPrompts() {
	for {
		select {
		case <-a.ctx.Done():
			return
		case prompt := <-a.cmdRunner.CodePromptChan():
			a.mu.Lock()
			a.activePrompt = prompt
			a.mu.Unlock()

			// Fire event to UI
			runtime.EventsEmit(a.ctx, "steamcmd:auth_required", prompt.PromptType)
		}
	}
}

func (a *App) checkSteamCmdOnStartup() {
	// Let Wails UI mount to WebView first
	time.Sleep(1500 * time.Millisecond)

	cfg := a.cfgManager.GetConfig()
	if _, err := os.Stat(cfg.SteamCmdPath); os.IsNotExist(err) {
		choice, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.QuestionDialog,
			Title:         "SteamCMD Not Found",
			Message:       "SteamCMD was not found at the configured path.\nWould you like to download and install it automatically?",
			Buttons:       []string{"Yes", "No"},
			DefaultButton: "Yes",
		})
		if err == nil && choice == "Yes" {
			logFn := func(event domain.LogEvent) {
				runtime.EventsEmit(a.ctx, "steamcmd:log", event)
			}
			ctx, cancel := context.WithTimeout(a.ctx, 5*time.Minute)
			defer cancel()
			
			err = steamcmd.VerifyOrInstall(ctx, cfg.SteamCmdPath, logFn)
			if err != nil {
				runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
					Type:    runtime.ErrorDialog,
					Title:   "Installation Failed",
					Message: fmt.Sprintf("Failed to install SteamCMD: %s", err.Error()),
				})
			}
		}
	}
}

// Bindings

func (a *App) GetConfig() domain.AppConfig {
	return a.cfgManager.GetConfig()
}

func (a *App) SaveConfig(cfg domain.AppConfig) error {
	// Verify if the path changed and SteamCMD does not exist at new path
	oldCfg := a.cfgManager.GetConfig()
	if oldCfg.SteamCmdPath != cfg.SteamCmdPath {
		if _, err := os.Stat(cfg.SteamCmdPath); os.IsNotExist(err) {
			choice, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:          runtime.QuestionDialog,
				Title:         "SteamCMD Path Warning",
				Message:       fmt.Sprintf("SteamCMD was not found at: %s\nWould you like to install it automatically at this location?", cfg.SteamCmdPath),
				Buttons:       []string{"Yes", "No"},
				DefaultButton: "Yes",
			})
			if err == nil && choice == "Yes" {
				go func() {
					logFn := func(event domain.LogEvent) {
						runtime.EventsEmit(a.ctx, "steamcmd:log", event)
					}
					ctx, cancel := context.WithTimeout(a.ctx, 5*time.Minute)
					defer cancel()
					_ = steamcmd.VerifyOrInstall(ctx, cfg.SteamCmdPath, logFn)
				}()
			}
		}
	}

	err := a.cfgManager.UpdateConfig(cfg)
	if err == nil {
		a.cmdRunner.UpdatePath(cfg.SteamCmdPath)
	}
	return err
}

func (a *App) CheckSteamCmd() (string, error) {
	cfg := a.cfgManager.GetConfig()
	if _, err := os.Stat(cfg.SteamCmdPath); os.IsNotExist(err) {
		return "not_found", nil
	}
	return "found", nil
}

func (a *App) ForceInstallSteamCmd() error {
	cfg := a.cfgManager.GetConfig()
	logFn := func(event domain.LogEvent) {
		runtime.EventsEmit(a.ctx, "steamcmd:log", event)
	}
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Minute)
	defer cancel()

	return steamcmd.VerifyOrInstall(ctx, cfg.SteamCmdPath, logFn)
}

func (a *App) AddDownload(workshopURL string, manualAppID string) (string, error) {
	workshopID, err := extractWorkshopID(workshopURL)
	if err != nil {
		return "", err
	}

	appID := manualAppID
	if appID == "" {
		appID = extractAppIDFromURL(workshopURL)
		if appID == "" {
			return "", fmt.Errorf("could not determine AppID. Please specify manually")
		}
	}

	jobID := a.downManager.AddJob(workshopID, appID)
	return jobID, nil
}

func (a *App) GetJobs() []domain.DownloadJob {
	return a.downManager.GetJobs()
}

func (a *App) CancelJob(jobID string) {
	a.downManager.CancelJob(jobID)
}

func (a *App) SubmitSteamCode(code string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.activePrompt == nil {
		return fmt.Errorf("no active steam guard prompt found")
	}

	// Send code back to runner and clean up
	select {
	case a.activePrompt.Response <- code:
	default:
	}
	a.activePrompt = nil
	return nil
}

func (a *App) CancelSteamCodePrompt() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.activePrompt == nil {
		return nil
	}

	// Send empty string to signal cancellation/close stdin
	select {
	case a.activePrompt.Response <- "":
	default:
	}
	a.activePrompt = nil
	return nil
}

// LoginSteam executes direct SteamCMD login to authenticate account and resolve 2FA
func (a *App) LoginSteam(username, password string) error {
	cfg := a.cfgManager.GetConfig()
	cfg.Username = username
	if err := a.cfgManager.UpdateConfig(cfg); err != nil {
		return fmt.Errorf("failed to save login configuration: %w", err)
	}

	logFn := func(event domain.LogEvent) {
		runtime.EventsEmit(a.ctx, "steamcmd:log", event)
	}

	// Verify or install SteamCMD dynamically
	if err := steamcmd.VerifyOrInstall(a.ctx, cfg.SteamCmdPath, logFn); err != nil {
		return fmt.Errorf("steamcmd verification/run failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Minute)
	defer cancel()

	err := a.cmdRunner.Login(ctx, username, password, logFn)
	if err != nil {
		return fmt.Errorf("steamcmd login failed: %w", err)
	}

	return nil
}

// Helpers

func extractWorkshopID(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		m, _ := regexp.MatchString(`^\d+$`, urlStr)
		if m {
			return urlStr, nil
		}
		return "", fmt.Errorf("invalid workshop URL or ID")
	}

	id := u.Query().Get("id")
	if id == "" {
		re := regexp.MustCompile(`sharedfiles/filedetails/\?id=(\d+)`)
		matches := re.FindStringSubmatch(urlStr)
		if len(matches) > 1 {
			return matches[1], nil
		}
		
		m, _ := regexp.MatchString(`^\d+$`, urlStr)
		if m {
			return urlStr, nil
		}

		return "", fmt.Errorf("could not locate workshop item 'id' parameter in URL")
	}
	return id, nil
}

func extractAppIDFromURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	appid := u.Query().Get("appid")
	if appid != "" {
		return appid
	}
	return ""
}
