package downloader

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"workshop-downloader/internal/config"
	"workshop-downloader/internal/domain"
	"workshop-downloader/internal/logger"
	"workshop-downloader/internal/steamcmd"
)

type EventBroadcaster interface {
	Emit(eventName string, data interface{})
}

type DownloadManager struct {
	mu           sync.RWMutex
	jobs         map[string]*domain.DownloadJob
	queue        []string // Job IDs in queue
	activeJobID  string
	runner       *steamcmd.SteamCmdRunner
	cfgManager   *config.ConfigManager
	broadcaster  EventBroadcaster
	cancelActive context.CancelFunc
	running      bool
	workerCtx    context.Context
	workerCancel context.CancelFunc
}

func NewDownloadManager(
	runner *steamcmd.SteamCmdRunner,
	cfgManager *config.ConfigManager,
	broadcaster EventBroadcaster,
) *DownloadManager {
	ctx, cancel := context.WithCancel(context.Background())
	dm := &DownloadManager{
		jobs:         make(map[string]*domain.DownloadJob),
		queue:        []string{},
		runner:       runner,
		cfgManager:   cfgManager,
		broadcaster:  broadcaster,
		workerCtx:    ctx,
		workerCancel: cancel,
	}

	return dm
}

func (dm *DownloadManager) Start() {
	dm.mu.Lock()
	if dm.running {
		dm.mu.Unlock()
		return
	}
	dm.running = true
	dm.mu.Unlock()

	go dm.processQueue()
}

func (dm *DownloadManager) Stop() {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	if !dm.running {
		return
	}
	dm.workerCancel()
	if dm.cancelActive != nil {
		dm.cancelActive()
	}
	dm.running = false
}

func (dm *DownloadManager) AddJob(workshopID, appID string) string {
	return dm.AddJobWithMeta(workshopID, appID, "", "")
}

func (dm *DownloadManager) AddJobWithMeta(workshopID, appID, title, previewURL string) string {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	jobID := fmt.Sprintf("job_%d", time.Now().UnixNano())
	job := &domain.DownloadJob{
		ID:         jobID,
		WorkshopID: workshopID,
		AppID:      appID,
		Title:      title,
		PreviewURL: previewURL,
		Status:     domain.StatusQueued,
		Progress:   0,
		CreatedAt:  time.Now(),
	}

	dm.jobs[jobID] = job
	dm.queue = append(dm.queue, jobID)

	logger.WriteLog(fmt.Sprintf("Queue: Added new download job ID=%s, WorkshopID=%s, AppID=%s, Title=%s", jobID, workshopID, appID, title))
	dm.broadcastJobChange(job)
	return jobID
}

func (dm *DownloadManager) GetJobs() []domain.DownloadJob {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	allJobs := []domain.DownloadJob{}
	for _, id := range dm.queue {
		if j, ok := dm.jobs[id]; ok {
			allJobs = append(allJobs, *j)
		}
	}
	for _, j := range dm.jobs {
		found := false
		for _, qid := range dm.queue {
			if qid == j.ID {
				found = true
				break
			}
		}
		if !found {
			allJobs = append(allJobs, *j)
		}
	}
	return allJobs
}

func (dm *DownloadManager) CancelJob(jobID string) {
	dm.mu.Lock()
	job, exists := dm.jobs[jobID]
	if !exists {
		dm.mu.Unlock()
		return
	}

	logger.WriteLog(fmt.Sprintf("Queue: Job cancellation requested ID=%s, CurrentStatus=%s", jobID, job.Status))
	if job.Status == domain.StatusRunning {
		if dm.cancelActive != nil {
			dm.cancelActive()
		}
	} else if job.Status == domain.StatusQueued {
		job.Status = domain.StatusCancelled
		now := time.Now()
		job.FinishedAt = &now
		dm.removeFromQueue(jobID)
		dm.broadcastJobChange(job)
	}
	dm.mu.Unlock()
}

func (dm *DownloadManager) RetryJob(jobID string) error {
	dm.mu.Lock()
	job, exists := dm.jobs[jobID]
	if !exists {
		dm.mu.Unlock()
		return fmt.Errorf("job not found")
	}

	if job.Status == domain.StatusRunning || job.Status == domain.StatusQueued {
		dm.mu.Unlock()
		return fmt.Errorf("job is already active or in queue")
	}

	// Reset job fields
	job.Status = domain.StatusQueued
	job.Progress = 0
	job.ErrorMsg = ""
	now := time.Now()
	job.CreatedAt = now
	job.StartedAt = nil
	job.FinishedAt = nil

	// Append back to queue list if not present
	found := false
	for _, id := range dm.queue {
		if id == jobID {
			found = true
			break
		}
	}
	if !found {
		dm.queue = append(dm.queue, jobID)
	}

	logger.WriteLog(fmt.Sprintf("Queue: Retrying job ID=%s, WorkshopID=%s, AppID=%s", jobID, job.WorkshopID, job.AppID))
	dm.broadcastJobChange(job)
	dm.mu.Unlock()
	return nil
}

func (dm *DownloadManager) DeleteJob(jobID string) {
	dm.mu.Lock()
	job, exists := dm.jobs[jobID]
	if !exists {
		dm.mu.Unlock()
		return
	}

	logger.WriteLog(fmt.Sprintf("Queue: Deleting job ID=%s", jobID))

	if job.Status == domain.StatusRunning {
		if dm.cancelActive != nil {
			dm.cancelActive()
		}
	}

	dm.removeFromQueue(jobID)
	delete(dm.jobs, jobID)
	dm.mu.Unlock()

	dm.broadcaster.Emit("job:deleted", jobID)
}

func (dm *DownloadManager) processQueue() {
	for {
		select {
		case <-dm.workerCtx.Done():
			return
		default:
		}

		var activeJob *domain.DownloadJob
		dm.mu.Lock()
		if len(dm.queue) > 0 {
			var nextJobIdx = -1
			for i, id := range dm.queue {
				if dm.jobs[id].Status == domain.StatusQueued {
					nextJobIdx = i
					break
				}
			}

			if nextJobIdx != -1 {
				id := dm.queue[nextJobIdx]
				activeJob = dm.jobs[id]
				dm.activeJobID = id
				activeJob.Status = domain.StatusRunning
				now := time.Now()
				activeJob.StartedAt = &now
				logger.WriteLog(fmt.Sprintf("Queue: Processing active job ID=%s, WorkshopID=%s, AppID=%s", id, activeJob.WorkshopID, activeJob.AppID))
				dm.broadcastJobChange(activeJob)
			}
		}
		dm.mu.Unlock()

		if activeJob == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		err := dm.runJob(activeJob)

		dm.mu.Lock()
		now := time.Now()
		activeJob.FinishedAt = &now
		if err != nil {
			if activeJob.Status != domain.StatusCancelled {
				activeJob.Status = domain.StatusFailed
				activeJob.ErrorMsg = err.Error()
				logger.WriteLog(fmt.Sprintf("Queue: Job ID=%s FAILED with error: %s", activeJob.ID, err.Error()))
			} else {
				logger.WriteLog(fmt.Sprintf("Queue: Job ID=%s was CANCELLED", activeJob.ID))
			}
		} else {
			activeJob.Status = domain.StatusSuccess
			activeJob.Progress = 100.0
			logger.WriteLog(fmt.Sprintf("Queue: Job ID=%s COMPLETED successfully", activeJob.ID))
		}
		dm.removeFromQueue(activeJob.ID)
		dm.activeJobID = ""
		dm.broadcastJobChange(activeJob)
		dm.mu.Unlock()
	}
}

func (dm *DownloadManager) runJob(job *domain.DownloadJob) error {
	ctx, cancel := context.WithCancel(dm.workerCtx)
	dm.mu.Lock()
	dm.cancelActive = cancel
	dm.mu.Unlock()
	defer cancel()

	cfg := dm.cfgManager.GetConfig()

	logFn := func(event domain.LogEvent) {
		event.JobID = job.ID
		dm.broadcaster.Emit("steamcmd:log", event)
	}

	progressFn := func(progress float64) {
		dm.mu.Lock()
		job.Progress = progress
		dm.mu.Unlock()

		dm.broadcaster.Emit("job:progress", domain.ProgressEvent{
			JobID:    job.ID,
			Progress: progress,
		})
	}

	// 1. Verify or install SteamCMD dynamically
	err := steamcmd.VerifyOrInstall(ctx, cfg.SteamCmdPath, logFn)
	if err != nil {
		return fmt.Errorf("steamcmd init failed: %w", err)
	}

	// 2. Perform the download operation
	err = dm.runner.DownloadItem(ctx, job.AppID, job.WorkshopID, cfg.Username, "", logFn, progressFn)
	if err != nil {
		return err
	}

	// 3. Move the downloaded item from SteamCMD cache to custom OutputDir
	steamCmdDir := filepath.Dir(cfg.SteamCmdPath)
	srcDir := filepath.Join(steamCmdDir, "steamapps", "workshop", "content", job.AppID, job.WorkshopID)
	destDir := filepath.Join(cfg.OutputDir, job.AppID, job.WorkshopID)

	logFn(domain.LogEvent{
		Message:   fmt.Sprintf("Moving downloaded files from cache to target directory (%s)...", destDir),
		Timestamp: time.Now(),
	})

	if _, statErr := os.Stat(srcDir); statErr == nil {
		if moveErr := moveFiles(srcDir, destDir); moveErr != nil {
			return fmt.Errorf("failed to move files to destination: %w", moveErr)
		}
		logFn(domain.LogEvent{
			Message:   "Files successfully moved to: " + destDir,
			Timestamp: time.Now(),
		})
	} else {
		return fmt.Errorf("download reported success, but content folder was not found at expected path: %s", srcDir)
	}

	return nil
}

func (dm *DownloadManager) removeFromQueue(jobID string) {
	newQueue := []string{}
	for _, id := range dm.queue {
		if id != jobID {
			newQueue = append(newQueue, id)
		}
	}
	dm.queue = newQueue
}

func (dm *DownloadManager) broadcastJobChange(job *domain.DownloadJob) {
	dm.broadcaster.Emit("job:status", *job)
}

func moveFiles(src, dest string) error {
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	files, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range files {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if err := os.Rename(srcPath, destPath); err != nil {
			if err := copyFile(srcPath, destPath); err != nil {
				return fmt.Errorf("failed to copy file: %w", err)
			}
			_ = os.Remove(srcPath)
		}
	}
	_ = os.Remove(src)
	return nil
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
