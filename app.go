package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"gorecorder/backend"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	recorder *backend.Recorder
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		recorder: backend.NewRecorder(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Start background status ticker to broadcast live status & duration to frontend
	go a.statusBroadcastLoop()
}

// shutdown is called when the app closes
func (a *App) shutdown(ctx context.Context) {
	if a.recorder != nil {
		a.recorder.Cleanup()
	}
}

func (a *App) statusBroadcastLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			if a.recorder != nil {
				st := a.recorder.GetStatus()
				wailsRuntime.EventsEmit(a.ctx, "recorder:status", st)
			}
		}
	}
}

// StartRecording starts recording with selected configuration
func (a *App) StartRecording(cfg backend.RecordConfig) error {
	return a.recorder.Start(cfg)
}

// PauseRecording pauses the current recording via SIGSTOP
func (a *App) PauseRecording() error {
	return a.recorder.Pause()
}

// ResumeRecording resumes recording via SIGCONT
func (a *App) ResumeRecording() error {
	return a.recorder.Resume()
}

// TogglePauseResume toggles pause/resume (for hotkey F9)
func (a *App) TogglePauseResume() (string, error) {
	st, err := a.recorder.TogglePauseResume()
	return string(st), err
}

// StopRecording finalizes recording, remuxes to MP4, and cleans temp files
func (a *App) StopRecording() (string, error) {
	return a.recorder.Stop()
}

// IsRecording returns if recording is in progress
func (a *App) IsRecording() bool {
	return a.recorder.IsRecording()
}

// GetCurrentDuration returns active recording duration in seconds
func (a *App) GetCurrentDuration() int64 {
	return a.recorder.CurrentDuration()
}

// GetStatus returns the current recording status
func (a *App) GetStatus() backend.RecordingStatus {
	return a.recorder.GetStatus()
}

// GetAudioDevices returns list of microphones and system audio monitor sinks
func (a *App) GetAudioDevices() struct {
	Mics     []backend.AudioDeviceInfo `json:"mics"`
	Monitors []backend.AudioDeviceInfo `json:"monitors"`
} {
	mics, monitors := backend.GetAudioDevices()
	return struct {
		Mics     []backend.AudioDeviceInfo `json:"mics"`
		Monitors []backend.AudioDeviceInfo `json:"monitors"`
	}{
		Mics:     mics,
		Monitors: monitors,
	}
}

// GetEncoders returns detected encoders
func (a *App) GetEncoders() []backend.EncoderInfo {
	return backend.GetAllAvailableEncoders()
}

// GetRecentRecordings returns list of recently completed recordings
func (a *App) GetRecentRecordings() []backend.RecordedFileEntry {
	return a.recorder.GetRecentRecordings()
}

// GetDefaultOutputDir returns the user's default Videos directory
func (a *App) GetDefaultOutputDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/tmp"
	}
	videosDir := filepath.Join(home, "Videos")
	_ = os.MkdirAll(videosDir, 0755)
	return videosDir
}

// SelectOutputDir opens a native folder selection dialog
func (a *App) SelectOutputDir() (string, error) {
	defaultDir := a.GetDefaultOutputDir()
	dir, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		DefaultDirectory: defaultDir,
		Title:            "Select Recording Output Directory",
	})
	if err != nil {
		return "", err
	}
	if dir == "" {
		return defaultDir, nil
	}
	return dir, nil
}

// OpenFolder opens the file manager highlighting the given directory or file
func (a *App) OpenFolder(targetPath string) error {
	if targetPath == "" {
		targetPath = a.GetDefaultOutputDir()
	}

	fi, err := os.Stat(targetPath)
	if err == nil && !fi.IsDir() {
		targetPath = filepath.Dir(targetPath)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("xdg-open", targetPath)
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", targetPath)
	} else {
		cmd = exec.Command("explorer", targetPath)
	}

	return cmd.Start()
}

// PlayVideo opens the recorded video in the default media player
func (a *App) PlayVideo(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("file does not exist: %s", filePath)
	}
	cmd := exec.Command("xdg-open", filePath)
	return cmd.Start()
}
