package backend

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// FFmpegProcess wraps the running recording process
type FFmpegProcess struct {
	cmd      *exec.Cmd
	cancelFn context.CancelFunc
	tempMKV  string
	finalMP4 string
	encoder  EncoderInfo
	isKilled bool
}

// BuildAndStartFFmpeg starts recording session via PipeWire / GStreamer / FFmpeg
func BuildAndStartFFmpeg(
	ctx context.Context,
	cfg RecordConfig,
	encoder EncoderInfo,
	nodeID uint32,
	pwFD int,
	pwFile *os.File,
) (*FFmpegProcess, error) {
	cmdCtx, cancel := context.WithCancel(ctx)

	outDir := cfg.OutputDir
	if outDir == "" {
		home, _ := os.UserHomeDir()
		outDir = filepath.Join(home, "Videos")
		_ = os.MkdirAll(outDir, 0755)
	}

	timestamp := time.Now().Format("2006-01-02-15-04")
	baseName := fmt.Sprintf("golang-lesson-%s", timestamp)
	finalMP4 := filepath.Join(outDir, fmt.Sprintf("%s.mp4", baseName))

	if _, err := os.Stat(finalMP4); err == nil {
		finalMP4 = filepath.Join(outDir, fmt.Sprintf("%s-%02d.mp4", baseName, time.Now().Second()))
	}

	tempMKV := filepath.Join(outDir, fmt.Sprintf("temp_recording_%d.mkv", time.Now().UnixNano()))

	fps := cfg.FPS
	if fps <= 0 {
		fps = 60
	}

	gstArgs := BuildGstreamerArgs(cfg, tempMKV, nodeID, pwFile != nil)

	cmd := exec.CommandContext(cmdCtx, "gst-launch-1.0", gstArgs...)
	if pwFile != nil {
		cmd.ExtraFiles = []*os.File{pwFile}
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to start recording process: %w", err)
	}

	return &FFmpegProcess{
		cmd:      cmd,
		cancelFn: cancel,
		tempMKV:  tempMKV,
		finalMP4: finalMP4,
		encoder:  encoder,
	}, nil
}

// Pause sends SIGSTOP to the process group
func (p *FFmpegProcess) Pause() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return fmt.Errorf("process not running")
	}
	pid := p.cmd.Process.Pid
	return syscall.Kill(-pid, syscall.SIGSTOP)
}

// Resume sends SIGCONT to the process group
func (p *FFmpegProcess) Resume() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return fmt.Errorf("process not running")
	}
	pid := p.cmd.Process.Pid
	return syscall.Kill(-pid, syscall.SIGCONT)
}

// Stop gracefully finalizes recording with SIGINT, waits for file flush, then remuxes MKV -> MP4
func (p *FFmpegProcess) Stop() (string, error) {
	if p.cmd == nil || p.cmd.Process == nil {
		return "", fmt.Errorf("process not running")
	}

	pid := p.cmd.Process.Pid

	// Resume before sending SIGINT in case it was paused
	_ = syscall.Kill(-pid, syscall.SIGCONT)
	_ = syscall.Kill(-pid, syscall.SIGINT)

	done := make(chan error, 1)
	go func() {
		done <- p.cmd.Wait()
	}()

	select {
	case <-done:
	case <-time.After(15 * time.Second):
		_ = syscall.Kill(-pid, syscall.SIGTERM)
		time.Sleep(1 * time.Second)
		_ = syscall.Kill(-pid, syscall.SIGKILL)
	}

	if p.cancelFn != nil {
		p.cancelFn()
	}

	// Check if temp MKV exists
	if _, err := os.Stat(p.tempMKV); err != nil {
		return "", fmt.Errorf("recording file was not generated: %w", err)
	}

	// Remux MKV -> MP4: ffmpeg -i temp_recording.mkv -codec copy final.mp4
	err := RemuxMKVToMP4(p.tempMKV, p.finalMP4)
	if err != nil {
		return p.tempMKV, fmt.Errorf("remux to MP4 failed, saved as MKV: %w", err)
	}

	// Delete temporary MKV after successful remux
	_ = os.Remove(p.tempMKV)

	return p.finalMP4, nil
}

// RemuxMKVToMP4 performs lossless stream copy remux from MKV to MP4
func RemuxMKVToMP4(mkvPath, mp4Path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-hide_banner", "-loglevel", "error", "-y",
		"-i", mkvPath,
		"-codec", "copy",
		"-movflags", "+faststart",
		mp4Path,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("remux error: %s (%w)", string(out), err)
	}
	return nil
}

// ForceKill ensures no zombie process is left on app exit
func (p *FFmpegProcess) ForceKill() {
	if p == nil || p.cmd == nil || p.cmd.Process == nil {
		return
	}
	pid := p.cmd.Process.Pid
	_ = syscall.Kill(-pid, syscall.SIGKILL)
	if p.cancelFn != nil {
		p.cancelFn()
	}
	if p.tempMKV != "" {
		_ = os.Remove(p.tempMKV)
	}
}

// BuildGstreamerArgs constructs the complete GStreamer argument slice for reliable recording
func BuildGstreamerArgs(cfg RecordConfig, tempMKV string, nodeID uint32, hasPWFile bool) []string {
	fps := cfg.FPS
	if fps <= 0 {
		fps = 60
	}

	var gstArgs []string
	gstArgs = append(gstArgs,
		"-e",
		"matroskamux", "name=mux", "!", "filesink", fmt.Sprintf("location=%s", tempMKV),
	)

	// Video pipeline
	if hasPWFile && nodeID > 0 {
		gstArgs = append(gstArgs,
			"pipewiresrc", "fd=3", fmt.Sprintf("path=%d", nodeID), "do-timestamp=true", "keepalive-time=1000",
		)
	} else if nodeID > 0 {
		gstArgs = append(gstArgs,
			"pipewiresrc", fmt.Sprintf("path=%d", nodeID), "do-timestamp=true", "keepalive-time=1000",
		)
	} else {
		// Fallback x11
		display := os.Getenv("DISPLAY")
		if display == "" {
			display = ":0"
		}
		gstArgs = append(gstArgs,
			"ximagesrc", fmt.Sprintf("display-name=%s", display), "use-damage=0",
		)
	}

	// Resolution scaling
	var scalePipeline []string
	switch cfg.Resolution {
	case Res480p:
		scalePipeline = []string{"!", "videoscale", "!", "video/x-raw,width=854,height=480"}
	case Res720p:
		scalePipeline = []string{"!", "videoscale", "!", "video/x-raw,width=1280,height=720"}
	case Res1080p:
		scalePipeline = []string{"!", "videoscale", "!", "video/x-raw,width=1920,height=1080"}
	case Res1440p:
		scalePipeline = []string{"!", "videoscale", "!", "video/x-raw,width=2560,height=1440"}
	}

	gstArgs = append(gstArgs, "!", "videoconvert")
	if len(scalePipeline) > 0 {
		gstArgs = append(gstArgs, scalePipeline...)
	}

	// Video encoder (Fast & Crisp with keyframes every 2s, parsed and queued with generous 10s buffer)
	gstArgs = append(gstArgs,
		"!", "videorate",
		"!", fmt.Sprintf("video/x-raw,framerate=%d/1", fps),
		"!", "x264enc", "speed-preset=veryfast", "tune=zerolatency", "bitrate=5000", fmt.Sprintf("key-int-max=%d", fps*2),
		"!", "h264parse",
		"!", "queue", "max-size-buffers=0", "max-size-time=10000000000", "max-size-bytes=0",
		"!", "mux.video_0",
	)

	// Audio pipeline (AAC 192k with audiorate sample correction, hard-resync, aacparse, and generous queues)
	micDev := cfg.MicDevice
	if micDev == "" {
		micDev = "default"
	}
	sysDev := cfg.SysDevice
	if sysDev == "" {
		sysDev = "@DEFAULT_SINK@.monitor"
	}

	switch cfg.AudioSource {
	case AudioSourceMic:
		gstArgs = append(gstArgs,
			"pulsesrc", "buffer-time=2000000", "latency-time=10000", fmt.Sprintf("device=%s", micDev),
			"!", "audioconvert",
			"!", "audioresample",
			"!", "audiorate",
			"!", "audio/x-raw,rate=48000,channels=2",
			"!", "avenc_aac", "bitrate=192000", "perfect-timestamp=true", "hard-resync=true",
			"!", "aacparse",
			"!", "queue", "max-size-buffers=0", "max-size-time=10000000000", "max-size-bytes=0",
			"!", "mux.audio_0",
		)
	case AudioSourceSystem:
		gstArgs = append(gstArgs,
			"pulsesrc", "buffer-time=2000000", "latency-time=10000", fmt.Sprintf("device=%s", sysDev),
			"!", "audioconvert",
			"!", "audioresample",
			"!", "audiorate",
			"!", "audio/x-raw,rate=48000,channels=2",
			"!", "avenc_aac", "bitrate=192000", "perfect-timestamp=true", "hard-resync=true",
			"!", "aacparse",
			"!", "queue", "max-size-buffers=0", "max-size-time=10000000000", "max-size-bytes=0",
			"!", "mux.audio_0",
		)
	case AudioSourceBoth:
		gstArgs = append(gstArgs,
			"audiomixer", "name=mix", "ignore-inactive-pads=true",
			"!", "audioconvert",
			"!", "audioresample",
			"!", "audiorate",
			"!", "audio/x-raw,rate=48000,channels=2",
			"!", "avenc_aac", "bitrate=192000", "perfect-timestamp=true", "hard-resync=true",
			"!", "aacparse",
			"!", "queue", "max-size-buffers=0", "max-size-time=10000000000", "max-size-bytes=0",
			"!", "mux.audio_0",
			"pulsesrc", "buffer-time=2000000", "latency-time=10000", fmt.Sprintf("device=%s", micDev),
			"!", "audioconvert", "!", "audioresample", "!", "audiorate", "!", "audio/x-raw,rate=48000,channels=2",
			"!", "queue", "max-size-buffers=0", "max-size-time=5000000000", "max-size-bytes=0", "leaky=downstream",
			"!", "mix.sink_0",
			"pulsesrc", "buffer-time=2000000", "latency-time=10000", fmt.Sprintf("device=%s", sysDev),
			"!", "audioconvert", "!", "audioresample", "!", "audiorate", "!", "audio/x-raw,rate=48000,channels=2",
			"!", "queue", "max-size-buffers=0", "max-size-time=5000000000", "max-size-bytes=0", "leaky=downstream",
			"!", "mix.sink_1",
		)
	}

	return gstArgs
}
