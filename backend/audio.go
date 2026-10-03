package backend

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
)

// AudioSourceResult contains the FFmpeg CLI arguments for capturing the requested audio sources
type AudioSourceResult struct {
	Args         []string
	AudioMapArgs []string
}

type pwDumpItem struct {
	Type string `json:"type"`
	Info struct {
		Props map[string]interface{} `json:"props"`
	} `json:"info"`
}

// GetAudioDevices returns all discovered microphone inputs and system monitor sinks
func GetAudioDevices() (mics []AudioDeviceInfo, monitors []AudioDeviceInfo) {
	// 1. Try native PipeWire pw-dump (available on all modern PipeWire systems)
	if pwOut, err := exec.Command("pw-dump").Output(); err == nil && len(pwOut) > 0 {
		var items []pwDumpItem
		if jsonErr := json.Unmarshal(pwOut, &items); jsonErr == nil {
			for _, item := range items {
				if item.Type != "PipeWire:Interface:Node" {
					continue
				}
				mediaClass, _ := item.Info.Props["media.class"].(string)
				nodeName, _ := item.Info.Props["node.name"].(string)
				nodeDesc, _ := item.Info.Props["node.description"].(string)
				nodeNick, _ := item.Info.Props["node.nick"].(string)

				if nodeName == "" {
					continue
				}
				label := nodeDesc
				if label == "" {
					label = nodeNick
				}
				if label == "" {
					label = simplifyAlsaName(nodeName)
				}

				if mediaClass == "Audio/Source" {
					mics = append(mics, AudioDeviceInfo{
						ID:          nodeName,
						Name:        nodeName,
						Description: "Microphone (" + label + ")",
						IsMonitor:   false,
					})
				} else if mediaClass == "Audio/Sink" {
					monitors = append(monitors, AudioDeviceInfo{
						ID:          nodeName + ".monitor",
						Name:        nodeName + ".monitor",
						Description: "System Audio (" + label + ")",
						IsMonitor:   true,
					})
				}
			}
		}
	}

	// 2. Fallback to pactl if pw-dump produced nothing
	if len(mics) == 0 && len(monitors) == 0 {
		cmd := exec.Command("pactl", "list", "sources", "short")
		out, err := cmd.Output()
		if err == nil {
			scanner := bufio.NewScanner(bytes.NewReader(out))
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					continue
				}
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					id := parts[0]
					name := parts[1]
					isMonitor := strings.HasSuffix(name, ".monitor")

					info := AudioDeviceInfo{
						ID:          id,
						Name:        name,
						Description: formatDeviceName(name),
						IsMonitor:   isMonitor,
					}

					if isMonitor {
						monitors = append(monitors, info)
					} else {
						mics = append(mics, info)
					}
				}
			}
		}
	}

	// 3. Add default entries if empty or as fallback options
	if len(mics) == 0 {
		mics = append(mics, AudioDeviceInfo{
			ID:          "default",
			Name:        "default",
			Description: "Default Microphone",
			IsDefault:   true,
			IsMonitor:   false,
		})
	}
	if len(monitors) == 0 {
		monitors = append(monitors, AudioDeviceInfo{
			ID:          "@DEFAULT_SINK@.monitor",
			Name:        "@DEFAULT_SINK@.monitor",
			Description: "Default System Audio Monitor",
			IsDefault:   true,
			IsMonitor:   true,
		})
	}

	return mics, monitors
}

func formatDeviceName(name string) string {
	if strings.Contains(name, "monitor") {
		clean := strings.TrimSuffix(name, ".monitor")
		return "System Audio (" + simplifyAlsaName(clean) + ")"
	}
	return "Microphone (" + simplifyAlsaName(name) + ")"
}

func simplifyAlsaName(name string) string {
	parts := strings.Split(name, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return name
}

// BuildAudioFFmpegArgs generates FFmpeg input flags and filters for chosen audio mode
// Audio bitrate requirement: AAC 192k
func BuildAudioFFmpegArgs(mode AudioSourceOption, micDev, sysDev string) AudioSourceResult {
	if mode == AudioSourceNone {
		return AudioSourceResult{}
	}

	if micDev == "" {
		micDev = "default"
	}
	if sysDev == "" {
		sysDev = "@DEFAULT_SINK@.monitor"
	}

	switch mode {
	case AudioSourceMic:
		return AudioSourceResult{
			Args: []string{
				"-f", "pulse",
				"-thread_queue_size", "1024",
				"-i", micDev,
			},
			AudioMapArgs: []string{
				"-c:a", "aac",
				"-b:a", "192k",
				"-ar", "48000",
			},
		}

	case AudioSourceSystem:
		return AudioSourceResult{
			Args: []string{
				"-f", "pulse",
				"-thread_queue_size", "1024",
				"-i", sysDev,
			},
			AudioMapArgs: []string{
				"-c:a", "aac",
				"-b:a", "192k",
				"-ar", "48000",
			},
		}

	case AudioSourceBoth:
		return AudioSourceResult{
			Args: []string{
				"-f", "pulse",
				"-thread_queue_size", "1024",
				"-i", micDev,
				"-f", "pulse",
				"-thread_queue_size", "1024",
				"-i", sysDev,
			},
			AudioMapArgs: []string{
				"-filter_complex", "[1:a][2:a]amix=inputs=2:duration=longest:dropout_transition=2[aout]",
				"-map", "0:v",
				"-map", "[aout]",
				"-c:a", "aac",
				"-b:a", "192k",
				"-ar", "48000",
			},
		}

	default:
		return AudioSourceResult{}
	}
}
