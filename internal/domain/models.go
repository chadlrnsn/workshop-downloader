package domain

import "time"

type JobStatus string

const (
	StatusQueued    JobStatus = "queued"
	StatusRunning   JobStatus = "running"
	StatusSuccess   JobStatus = "success"
	StatusFailed    JobStatus = "failed"
	StatusCancelled JobStatus = "cancelled"
)

type DownloadJob struct {
	ID         string     `json:"id"`
	WorkshopID string     `json:"workshopId"`
	AppID      string     `json:"appId"`
	Title      string     `json:"title,omitempty"`
	PreviewURL string     `json:"previewUrl,omitempty"`
	Status     JobStatus  `json:"status"`
	Progress   float64    `json:"progress"` // 0.0 to 100.0
	ErrorMsg   string     `json:"errorMsg,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	StartedAt  *time.Time `json:"startedAt,omitempty"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
}

type WorkshopItem struct {
	ID          string `json:"id"`
	AppID       string `json:"appId"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type SteamAccount struct {
	Username string `json:"username"`
	Remember bool   `json:"remember"`
}

type HistoryItem struct {
	ID           string    `json:"id"`
	WorkshopID   string    `json:"workshopId"`
	AppID        string    `json:"appId"`
	Title        string    `json:"title,omitempty"`
	PreviewURL   string    `json:"previewUrl,omitempty"`
	DownloadedAt time.Time `json:"downloadedAt"`
	FolderExists bool      `json:"folderExists"`
	Path         string    `json:"path,omitempty"`
}

type AppConfig struct {
	SteamCmdPath string `json:"steamCmdPath"`
	OutputDir    string `json:"outputDir"`
	AutoUpdate   bool   `json:"autoUpdate"`
	Username     string `json:"username"`
}

type LogEvent struct {
	JobID     string    `json:"jobId,omitempty"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	IsError   bool      `json:"isError"`
}

type ProgressEvent struct {
	JobID    string  `json:"jobId"`
	Progress float64 `json:"progress"`
	Speed    string  `json:"speed,omitempty"`
	ETA      string  `json:"eta,omitempty"`
}
