package backend

import (
	"testing"
)

func TestDetectBestEncoder(t *testing.T) {
	enc := DetectBestEncoder()
	t.Logf("Detected Best Encoder: %+v", enc)
	if enc.Type == "" {
		t.Fatalf("Expected valid encoder, got empty")
	}
}

func TestGetAudioDevices(t *testing.T) {
	mics, monitors := GetAudioDevices()
	t.Logf("Found %d mics and %d monitors", len(mics), len(monitors))
	if len(mics) == 0 {
		t.Errorf("Expected at least 1 microphone or fallback")
	}
	if len(monitors) == 0 {
		t.Errorf("Expected at least 1 monitor or fallback")
	}
}

func TestFormatSeconds(t *testing.T) {
	cases := []struct {
		secs     int64
		expected string
	}{
		{0, "00:00:00"},
		{59, "00:00:59"},
		{65, "00:01:05"},
		{3661, "01:01:01"},
	}

	for _, c := range cases {
		got := FormatSeconds(c.secs)
		if got != c.expected {
			t.Errorf("FormatSeconds(%d) = %s, expected %s", c.secs, got, c.expected)
		}
	}
}

func TestBuildAudioFFmpegArgs(t *testing.T) {
	none := BuildAudioFFmpegArgs(AudioSourceNone, "", "")
	if len(none.Args) != 0 {
		t.Errorf("Expected empty args for none audio")
	}

	mic := BuildAudioFFmpegArgs(AudioSourceMic, "default", "")
	if len(mic.Args) == 0 {
		t.Errorf("Expected non-empty args for mic audio")
	}

	both := BuildAudioFFmpegArgs(AudioSourceBoth, "default", "@DEFAULT_SINK@.monitor")
	if len(both.Args) == 0 || len(both.AudioMapArgs) == 0 {
		t.Errorf("Expected non-empty args and audio map for both")
	}
}

func TestBuildGstreamerArgs(t *testing.T) {
	cfgBoth := RecordConfig{
		Resolution:  Res1080p,
		FPS:         60,
		AudioSource: AudioSourceBoth,
	}

	args := BuildGstreamerArgs(cfgBoth, "/tmp/test.mkv", 42, true)
	argsStr := ""
	for _, a := range args {
		argsStr += a + " "
	}

	expectedElements := []string{
		"videorate",
		"x264enc",
		"h264parse",
		"audiomixer",
		"ignore-inactive-pads=true",
		"audiorate",
		"avenc_aac",
		"perfect-timestamp=true",
		"hard-resync=true",
		"aacparse",
		"max-size-time=10000000000",
		"buffer-time=2000000",
	}

	for _, elem := range expectedElements {
		found := false
		for _, a := range args {
			if a == elem {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected pipeline to contain %q, but was not found in: %s", elem, argsStr)
		}
	}
}
