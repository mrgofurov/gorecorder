package backend

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// DetectBestEncoder probes system for available encoders adhering to priority:
// 1. h264_vaapi
// 2. h264_nvenc
// 3. libx264
func DetectBestEncoder() EncoderInfo {
	// 1. Check h264_vaapi
	if vaapiDevice := findVAAPIDevice(); vaapiDevice != "" {
		if testEncoder(EncoderVAAPI, vaapiDevice) {
			return EncoderInfo{
				Type:        EncoderVAAPI,
				Name:        "Intel/AMD VAAPI (Hardware)",
				Description: "High-performance low-CPU hardware acceleration via VA-API",
				DevicePath:  vaapiDevice,
				IsHardware:  true,
			}
		}
	}

	// 2. Check h264_nvenc
	if testEncoder(EncoderNVENC, "") {
		return EncoderInfo{
			Type:        EncoderNVENC,
			Name:        "NVIDIA NVENC (Hardware)",
			Description: "High-speed NVIDIA GPU hardware encoder",
			IsHardware:  true,
		}
	}

	// 3. Fallback to libx264
	return EncoderInfo{
		Type:        EncoderLibx264,
		Name:        "libx264 (Software CPU)",
		Description: "Highly compatible CPU H.264 encoder (Fast preset)",
		IsHardware:  false,
	}
}

// GetAllAvailableEncoders returns all supported encoders found on this machine
func GetAllAvailableEncoders() []EncoderInfo {
	var list []EncoderInfo

	if dev := findVAAPIDevice(); dev != "" && testEncoder(EncoderVAAPI, dev) {
		list = append(list, EncoderInfo{
			Type:        EncoderVAAPI,
			Name:        "VAAPI Hardware (Intel/AMD)",
			Description: fmt.Sprintf("Hardware acceleration device: %s", dev),
			DevicePath:  dev,
			IsHardware:  true,
		})
	}

	if testEncoder(EncoderNVENC, "") {
		list = append(list, EncoderInfo{
			Type:        EncoderNVENC,
			Name:        "NVIDIA NVENC Hardware",
			Description: "Hardware acceleration via CUDA/NVENC",
			IsHardware:  true,
		})
	}

	list = append(list, EncoderInfo{
		Type:        EncoderLibx264,
		Name:        "libx264 CPU",
		Description: "CPU-based H.264 encoder (Fast/Efficient)",
		IsHardware:  false,
	})

	return list
}

func findVAAPIDevice() string {
	candidates := []string{
		"/dev/dri/renderD128",
		"/dev/dri/renderD129",
		"/dev/dri/card0",
		"/dev/dri/card1",
	}

	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			return c
		}
	}
	return ""
}

// testEncoder runs a 0.1-second dry test to verify FFmpeg can successfully init encoder
func testEncoder(encType EncoderType, devicePath string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var args []string
	switch encType {
	case EncoderVAAPI:
		if devicePath == "" {
			devicePath = "/dev/dri/renderD128"
		}
		args = []string{
			"-hide_banner", "-loglevel", "error",
			"-vaapi_device", devicePath,
			"-f", "lavfi", "-i", "testsrc=duration=0.1:size=64x64:rate=30",
			"-vf", "format=nv12,hwupload",
			"-c:v", "h264_vaapi",
			"-f", "null", "-",
		}
	case EncoderNVENC:
		args = []string{
			"-hide_banner", "-loglevel", "error",
			"-f", "lavfi", "-i", "testsrc=duration=0.1:size=64x64:rate=30",
			"-c:v", "h264_nvenc",
			"-f", "null", "-",
		}
	case EncoderLibx264:
		args = []string{
			"-hide_banner", "-loglevel", "error",
			"-f", "lavfi", "-i", "testsrc=duration=0.1:size=64x64:rate=30",
			"-c:v", "libx264", "-preset", "ultrafast",
			"-f", "null", "-",
		}
	default:
		return false
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		_ = out
		return false
	}
	return !strings.Contains(strings.ToLower(string(out)), "error")
}
