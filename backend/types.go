package backend

import "time"

// RecordingState defines the lifecycle state of the recorder
type RecordingState string

const (
	StateIdle       RecordingState = "IDLE"
	StateRecording  RecordingState = "REC"
	StatePaused     RecordingState = "PAUSED"
	StateProcessing RecordingState = "PROCESSING"
	StateError      RecordingState = "ERROR"
)

// ResolutionOption represents supported video resolutions
type ResolutionOption string

const (
	ResNative ResolutionOption = "native"
	Res480p   ResolutionOption = "480p"
	Res720p   ResolutionOption = "720p"
	Res1080p  ResolutionOption = "1080p"
	Res1440p  ResolutionOption = "1440p"
)

// AudioSourceOption represents audio recording choices
type AudioSourceOption string

const (
	AudioSourceNone       AudioSourceOption = "none"
	AudioSourceMic        AudioSourceOption = "mic"
	AudioSourceSystem     AudioSourceOption = "system"
	AudioSourceBoth       AudioSourceOption = "both"
)

// RecordConfig holds user settings for the screen recording session
type RecordConfig struct {
	Resolution  ResolutionOption  `json:"resolution"`
	FPS         int               `json:"fps"`
	AudioSource AudioSourceOption `json:"audioSource"`
	OutputDir   string            `json:"outputDir"`
	MicDevice   string            `json:"micDevice"`   // Optional specific mic source name
	SysDevice   string            `json:"sysDevice"`   // Optional specific system monitor source name
}

// EncoderType represents detected hardware / software video encoders
type EncoderType string

const (
	EncoderVAAPI  EncoderType = "h264_vaapi"
	EncoderNVENC  EncoderType = "h264_nvenc"
	EncoderLibx264 EncoderType = "libx264"
)

// EncoderInfo details encoder capabilities
type EncoderInfo struct {
	Type        EncoderType `json:"type"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	DevicePath  string      `json:"devicePath,omitempty"`
	IsHardware  bool        `json:"isHardware"`
}

// AudioDeviceInfo details an input/sink device
type AudioDeviceInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDefault   bool   `json:"isDefault"`
	IsMonitor   bool   `json:"isMonitor"`
}

// RecordingStatus represents live state communicated to frontend
type RecordingStatus struct {
	State        RecordingState `json:"state"`
	DurationSecs int64          `json:"durationSecs"`
	Formatted    string         `json:"formatted"`
	OutputFile   string         `json:"outputFile,omitempty"`
	TempFile     string         `json:"tempFile,omitempty"`
	EncoderUsed  string         `json:"encoderUsed,omitempty"`
	ErrorMessage string         `json:"errorMessage,omitempty"`
}

// RecordedFileEntry records metadata about recently created recordings
type RecordedFileEntry struct {
	Filename    string    `json:"filename"`
	FilePath    string    `json:"filePath"`
	SizeHuman   string    `json:"sizeHuman"`
	DurationSec int64     `json:"durationSec"`
	CreatedAt   time.Time `json:"createdAt"`
}
