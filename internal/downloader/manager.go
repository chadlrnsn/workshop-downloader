package downloader

import (
	"context"
	"fmt"
	"sync"
	"time"

	"workshop-downloader/internal/config"
	"workshop-downloader/internal/domain"
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
	dm.mu.Lock()
	defer dm.mu.Unlock()

	jobID := fmt.Sprintf("job_%d", time.Now().UnixNano())
	job := &domain.DownloadJob{
		ID:         jobID,
		WorkshopID: workshopID,
		AppID:      appID,
		Status:     domain.StatusQueued,
		Progress:   0,
		CreatedAt:  time.Now(),
	}

	dm.jobs[jobID] = job
	dm.queue = append(dm.queue, jobID)

	dm.broadcastJobChange(job)
	return jobID
}

func (dm *DownloadManager) GetJobs() []domain.DownloadJob {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	// Return jobs in order of creation/queue
	allJobs := []domain.DownloadJob{}
	for _, id := range dm.queue {
		if j, ok := dm.jobs[id]; ok {
			allJobs = append(allJobs, *j)
		}
	}
	// Add historical jobs that are not in queue if any
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
			// Find first queued job
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
				dm.broadcastJobChange(activeJob)
			}
		}
		dm.mu.Unlock()

		if activeJob == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// Execute job
		err := dm.runJob(activeJob)

		dm.mu.Lock()
		now := time.Now()
		activeJob.FinishedAt = &now
		if err != nil {
			if activeJob.Status != domain.StatusCancelled {
				activeJob.Status = domain.StatusFailed
				activeJob.ErrorMsg = err.Error()
			}
		} else {
			activeJob.Status = domain.StatusSuccess
			activeJob.Progress = 100.0
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

	// Broadcaster log utility
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
	return err
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
