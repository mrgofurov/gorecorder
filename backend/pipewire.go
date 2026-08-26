package backend

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
)

// PipeWireSession manages Wayland ScreenCast portal session and PipeWire stream node
type PipeWireSession struct {
	mu           sync.Mutex
	conn         *dbus.Conn
	sessionToken string
	sessionPath  dbus.ObjectPath
	nodeID       uint32
	pipewireFD   int
	pwFile       *os.File
	isWayland    bool
	closed       bool
}

// NewPipeWireSession detects desktop session type and prepares capture handles
func NewPipeWireSession() (*PipeWireSession, error) {
	isWayland := os.Getenv("XDG_SESSION_TYPE") == "wayland" || os.Getenv("WAYLAND_DISPLAY") != ""
	return &PipeWireSession{
		isWayland:  isWayland,
		pipewireFD: -1,
	}, nil
}

// RequestScreenCastSession initiates D-Bus screencast handshake via XDG Desktop Portal if on Wayland
func (s *PipeWireSession) RequestScreenCastSession(ctx context.Context) (uint32, int, *os.File, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conn, err := dbus.SessionBus()
	if err != nil {
		return 0, -1, nil, fmt.Errorf("failed to connect to session bus: %w", err)
	}
	s.conn = conn

	randNum := rand.New(rand.NewSource(time.Now().UnixNano())).Uint32()
	s.sessionToken = fmt.Sprintf("gorecorder_%d", randNum)
	rawSender := conn.Names()[0]
	sender := strings.ReplaceAll(strings.TrimPrefix(rawSender, ":"), ".", "_")

	portal := conn.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop")

	createToken := fmt.Sprintf("create_%d", randNum)
	createReqPath := dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/portal/desktop/request/%s/%s", sender, createToken))

	ch := make(chan *dbus.Signal, 25)
	conn.Signal(ch)
	_ = conn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.portal.Request"),
		dbus.WithMatchMember("Response"),
	)

	// 1. CreateSession
	createOpts := map[string]dbus.Variant{
		"session_handle_token": dbus.MakeVariant(s.sessionToken),
		"handle_token":         dbus.MakeVariant(createToken),
	}

	var sessReqPath dbus.ObjectPath
	err = portal.CallWithContext(ctx, "org.freedesktop.portal.ScreenCast.CreateSession", 0, createOpts).Store(&sessReqPath)
	if err != nil {
		return 0, -1, nil, fmt.Errorf("CreateSession DBus call failed: %w", err)
	}

	for {
		select {
		case sig := <-ch:
			if sig.Path == createReqPath {
				if len(sig.Body) >= 2 {
					if resMap, ok := sig.Body[1].(map[string]dbus.Variant); ok {
						if h, found := resMap["session_handle"]; found {
							s.sessionPath = dbus.ObjectPath(h.Value().(string))
						}
					}
				}
				goto Created
			}
		case <-time.After(3 * time.Second):
			s.sessionPath = dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/portal/desktop/session/%s/%s", sender, s.sessionToken))
			goto Created
		case <-ctx.Done():
			return 0, -1, nil, ctx.Err()
		}
	}

Created:
	// 2. SelectSources
	selectToken := fmt.Sprintf("select_%d", randNum)
	selectReqPath := dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/portal/desktop/request/%s/%s", sender, selectToken))

	selectOpts := map[string]dbus.Variant{
		"handle_token": dbus.MakeVariant(selectToken),
		"types":        dbus.MakeVariant(uint32(1 | 2)), // 1=monitor, 2=window
		"multiple":     dbus.MakeVariant(false),
		"cursor_mode":  dbus.MakeVariant(uint32(2)), // Embedded cursor
	}

	var selectCallPath dbus.ObjectPath
	_ = portal.CallWithContext(ctx, "org.freedesktop.portal.ScreenCast.SelectSources", 0, s.sessionPath, selectOpts).Store(&selectCallPath)

	for {
		select {
		case sig := <-ch:
			if sig.Path == selectReqPath {
				goto Selected
			}
		case <-time.After(3 * time.Second):
			goto Selected
		case <-ctx.Done():
			return 0, -1, nil, ctx.Err()
		}
	}

Selected:
	// 3. Start
	startToken := fmt.Sprintf("start_%d", randNum)
	startReqPath := dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/portal/desktop/request/%s/%s", sender, startToken))

	startOpts := map[string]dbus.Variant{
		"handle_token": dbus.MakeVariant(startToken),
	}

	var startCallPath dbus.ObjectPath
	_ = portal.CallWithContext(ctx, "org.freedesktop.portal.ScreenCast.Start", 0, s.sessionPath, "", startOpts).Store(&startCallPath)

	for {
		select {
		case sig := <-ch:
			if sig.Path == startReqPath {
				if len(sig.Body) >= 2 {
					if resMap, ok := sig.Body[1].(map[string]dbus.Variant); ok {
						if streamsVar, found := resMap["streams"]; found {
							if slice, ok := streamsVar.Value().([][]interface{}); ok && len(slice) > 0 {
								if nid, ok := slice[0][0].(uint32); ok {
									s.nodeID = nid
								}
							} else if arr, ok := streamsVar.Value().([]interface{}); ok && len(arr) > 0 {
								if tuple, ok := arr[0].([]interface{}); ok && len(tuple) > 0 {
									if nid, ok := tuple[0].(uint32); ok {
										s.nodeID = nid
									}
								}
							}
						}
					}
				}
				goto Started
			}
		case <-time.After(20 * time.Second):
			goto Started
		case <-ctx.Done():
			return 0, -1, nil, ctx.Err()
		}
	}

Started:
	// 4. OpenPipeWireRemote
	var pwFD dbus.UnixFD
	err = portal.CallWithContext(ctx, "org.freedesktop.portal.ScreenCast.OpenPipeWireRemote", 0, s.sessionPath, map[string]dbus.Variant{}).Store(&pwFD)
	if err == nil && int(pwFD) > 0 {
		s.pipewireFD = int(pwFD)
		s.pwFile = os.NewFile(uintptr(pwFD), "pipewire-remote")
	}

	return s.nodeID, s.pipewireFD, s.pwFile, nil
}

// Close gracefully closes the D-Bus screencast session
func (s *PipeWireSession) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}
	s.closed = true

	if s.pwFile != nil {
		_ = s.pwFile.Close()
		s.pwFile = nil
	}

	if s.conn != nil && s.sessionPath != "" {
		sessionObj := s.conn.Object("org.freedesktop.portal.Desktop", s.sessionPath)
		_ = sessionObj.Call("org.freedesktop.portal.Session.Close", 0)
		_ = s.conn.Close()
	}
}

// GetResolutionScaleFilter returns FFmpeg video filter for selected resolution
func GetResolutionScaleFilter(res ResolutionOption, isVAAPI bool) string {
	switch res {
	case Res480p:
		if isVAAPI {
			return "scale_vaapi=w=854:h=480"
		}
		return "scale=854:480:flags=bicubic"
	case Res720p:
		if isVAAPI {
			return "scale_vaapi=w=1280:h=720"
		}
		return "scale=1280:720:flags=bicubic"
	case Res1080p:
		if isVAAPI {
			return "scale_vaapi=w=1920:h=1080"
		}
		return "scale=1920:1080:flags=bicubic"
	case Res1440p:
		if isVAAPI {
			return "scale_vaapi=w=2560:h=1440"
		}
		return "scale=2560:1440:flags=bicubic"
	default:
		return ""
	}
}
