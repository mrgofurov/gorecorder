package backend

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Recorder manages the recording lifecycle, timer precision, and child process safety
type Recorder struct {
	mu           sync.RWMutex
	state        RecordingState
	config       RecordConfig
	process      *FFmpegProcess
	pwSession    *PipeWireSession
	bestEncoder  EncoderInfo
	startTime    time.Time
	pauseTime    time.Time
	totalPaused  time.Duration
	ctx          context.Context
	cancel       context.CancelFunc
	recentFiles  []RecordedFileEntry
	lastError    string
	statusTicker *time.Ticker
}

// NewRecorder creates and initializes the Recorder service
func NewRecorder() *Recorder {
	bestEnc := DetectBestEncoder()
	return &Recorder{
		state:       StateIdle,
		bestEncoder: bestEnc,
		recentFiles: make([]RecordedFileEntry, 0),
	}
}

// Start begins a new screen recording session
func (r *Recorder) Start(cfg RecordConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state == StateRecording || r.state == StatePaused {
		return fmt.Errorf("recording is already in progress")
	}

	r.state = StateProcessing
	r.config = cfg
	r.lastError = ""
	r.totalPaused = 0

	// Auto-detect best encoder
	r.bestEncoder = DetectBestEncoder()

	r.ctx, r.cancel = context.WithCancel(context.Background())

	// Initialize PipeWire session for Wayland
	var nodeID uint32
	var pwFD int = -1
	var pwFile *os.File

	pwSess, err := NewPipeWireSession()
	if err == nil {
		r.pwSession = pwSess
		nID, pFD, pFile, pwErr := pwSess.RequestScreenCastSession(r.ctx)
		if pwErr == nil {
			nodeID = nID
			pwFD = pFD
			pwFile = pFile
		} else {
			fmt.Printf("PipeWire screencast notice: %v\n", pwErr)
		}
	}

	// Launch recording process
	proc, err := BuildAndStartFFmpeg(r.ctx, cfg, r.bestEncoder, nodeID, pwFD, pwFile)
	if err != nil {
		r.state = StateError
		r.lastError = err.Error()
		if r.cancel != nil {
			r.cancel()
		}
		if r.pwSession != nil {
			r.pwSession.Close()
			r.pwSession = nil
		}
		return fmt.Errorf("failed to start recording: %w", err)
	}

	r.process = proc
	r.startTime = time.Now()
	r.state = StateRecording

	return nil
}

// Pause pauses the active recording via SIGSTOP without splitting files
func (r *Recorder) Pause() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != StateRecording {
		return fmt.Errorf("cannot pause: not currently recording")
	}

	if r.process == nil {
		return fmt.Errorf("no active recording process")
	}

	if err := r.process.Pause(); err != nil {
		return fmt.Errorf("failed to pause process: %w", err)
	}

	r.pauseTime = time.Now()
	r.state = StatePaused
	return nil
}

// Resume resumes the paused recording via SIGCONT
func (r *Recorder) Resume() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != StatePaused {
		return fmt.Errorf("cannot resume: recording is not paused")
	}

	if r.process == nil {
		return fmt.Errorf("no active recording process")
	}

	if err := r.process.Resume(); err != nil {
		return fmt.Errorf("failed to resume process: %w", err)
	}

	r.totalPaused += time.Since(r.pauseTime)
	r.state = StateRecording
	return nil
}

// TogglePauseResume toggles between recording and paused states (e.g. for F9 hotkey)
func (r *Recorder) TogglePauseResume() (RecordingState, error) {
	r.mu.RLock()
	st := r.state
	r.mu.RUnlock()

	if st == StateRecording {
		err := r.Pause()
		return StatePaused, err
	} else if st == StatePaused {
		err := r.Resume()
		return StateRecording, err
	}
	return st, fmt.Errorf("recording is not active")
}

// Stop gracefully finalizes the recording with SIGINT, converts MKV to MP4, and cleans temp files
func (r *Recorder) Stop() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != StateRecording && r.state != StatePaused {
		return "", fmt.Errorf("no active recording to stop")
	}

	r.state = StateProcessing
	finalDurationSecs := r.calculateDurationSecs()

	if r.process == nil {
		r.state = StateIdle
		if r.pwSession != nil {
			r.pwSession.Close()
			r.pwSession = nil
		}
		return "", fmt.Errorf("recording process reference was nil")
	}

	finalFile, err := r.process.Stop()
	r.process = nil

	if r.pwSession != nil {
		r.pwSession.Close()
		r.pwSession = nil
	}

	if r.cancel != nil {
		r.cancel()
		r.cancel = nil
	}

	if err != nil {
		r.state = StateError
		r.lastError = err.Error()
		return finalFile, err
	}

	r.state = StateIdle

	// Record metadata for UI history
	fi, statErr := os.Stat(finalFile)
	var sizeHuman string
	if statErr == nil {
		sizeHuman = formatBytes(fi.Size())
	}

	entry := RecordedFileEntry{
		Filename:    filepath.Base(finalFile),
		FilePath:    finalFile,
		SizeHuman:   sizeHuman,
		DurationSec: finalDurationSecs,
		CreatedAt:   time.Now(),
	}
	r.recentFiles = append([]RecordedFileEntry{entry}, r.recentFiles...)
	if len(r.recentFiles) > 20 {
		r.recentFiles = r.recentFiles[:20]
	}

	return finalFile, nil
}

// IsRecording returns true if recording is in progress
func (r *Recorder) IsRecording() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.state == StateRecording || r.state == StatePaused
}

// CurrentDuration returns duration in seconds excluding paused periods
func (r *Recorder) CurrentDuration() int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.calculateDurationSecs()
}

// GetStatus returns the current snapshot of recorder state and metrics
func (r *Recorder) GetStatus() RecordingStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dur := r.calculateDurationSecs()
	return RecordingStatus{
		State:        r.state,
		DurationSecs: dur,
		Formatted:    FormatSeconds(dur),
		EncoderUsed:  r.bestEncoder.Name,
		ErrorMessage: r.lastError,
	}
}

// GetRecentRecordings returns list of completed recordings
func (r *Recorder) GetRecentRecordings() []RecordedFileEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]RecordedFileEntry(nil), r.recentFiles...)
}

// Cleanup ensures any remaining background child processes are killed on application exit
func (r *Recorder) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.process != nil {
		r.process.ForceKill()
		r.process = nil
	}
	if r.pwSession != nil {
		r.pwSession.Close()
		r.pwSession = nil
	}
	if r.cancel != nil {
		r.cancel()
		r.cancel = nil
	}
	r.state = StateIdle
}

func (r *Recorder) calculateDurationSecs() int64 {
	if r.state == StateIdle || r.startTime.IsZero() {
		return 0
	}

	var elapsed time.Duration
	if r.state == StatePaused {
		elapsed = r.pauseTime.Sub(r.startTime) - r.totalPaused
	} else {
		elapsed = time.Since(r.startTime) - r.totalPaused
	}

	if elapsed < 0 {
		return 0
	}
	return int64(elapsed.Seconds())
}

// FormatSeconds formats seconds into HH:MM:SS
func FormatSeconds(secs int64) string {
	h := secs / 3600
	m := (secs % 3600) / 60
	s := secs % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
