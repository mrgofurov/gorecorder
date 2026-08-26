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
