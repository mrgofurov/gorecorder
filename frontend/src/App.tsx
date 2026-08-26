import React, { useEffect, useState, useCallback } from 'react';
import confetti from 'canvas-confetti';
import { Header } from './components/Header';
import { StatusIndicator } from './components/StatusIndicator';
import { ControlPanel } from './components/ControlPanel';
import { SettingsPanel } from './components/SettingsPanel';
import { RecentRecordings } from './components/RecentRecordings';
import { HotkeyBar } from './components/HotkeyBar';
import { Bridge } from './services/wailsBridge';
import { AudioDeviceInfo, EncoderInfo, RecordConfig, RecordedFileEntry, RecordingStatus } from './types';

export const App: React.FC = () => {
  const [config, setConfig] = useState<RecordConfig>({
    resolution: 'native',
    fps: 60,
    audioSource: 'both',
    outputDir: '',
  });

  const [status, setStatus] = useState<RecordingStatus>({
    state: 'IDLE',
    durationSecs: 0,
    formatted: '00:00:00',
  });

  const [encoders, setEncoders] = useState<EncoderInfo[]>([]);
  const [activeEncoder, setActiveEncoder] = useState<EncoderInfo | null>(null);
  const [audioDevices, setAudioDevices] = useState<{ mics: AudioDeviceInfo[]; monitors: AudioDeviceInfo[] }>({
    mics: [],
    monitors: [],
  });
  const [recentFiles, setRecentFiles] = useState<RecordedFileEntry[]>([]);
  const [notification, setNotification] = useState<{ message: string; type: 'success' | 'info' | 'error' } | null>(null);

  // Initialize data on mount
  useEffect(() => {
    const initData = async () => {
      try {
        const [defaultDir, encList, audioList, historyList] = await Promise.all([
          Bridge.getDefaultOutputDir(),
          Bridge.getEncoders(),
          Bridge.getAudioDevices(),
          Bridge.getRecentRecordings(),
        ]);

        setConfig((prev) => ({ ...prev, outputDir: defaultDir }));
        setEncoders(encList);
        if (encList.length > 0) {
          setActiveEncoder(encList[0]);
        }
        setAudioDevices(audioList);
        setRecentFiles(historyList);
      } catch (err) {
        console.error('Failed to load initial recorder configuration:', err);
      }
    };

    initData();

    // Subscribe to status events from Go backend
    const unsubscribe = Bridge.onStatusChange((newStatus) => {
      setStatus(newStatus);
    });

    return () => {
      unsubscribe();
    };
  }, []);

  const showToast = (message: string, type: 'success' | 'info' | 'error' = 'info') => {
    setNotification({ message, type });
    setTimeout(() => setNotification(null), 4000);
  };

  // Actions
  const handleStart = async () => {
    try {
      await Bridge.startRecording(config);
      setStatus((prev) => ({ ...prev, state: 'REC', durationSecs: 0, formatted: '00:00:00', errorMessage: undefined }));
      showToast('Recording started smoothly', 'info');
    } catch (err: any) {
      showToast(`Start failed: ${err?.message || err}`, 'error');
      setStatus((prev) => ({ ...prev, state: 'ERROR', errorMessage: String(err) }));
    }
  };

  const handlePause = async () => {
    try {
      await Bridge.pauseRecording();
      setStatus((prev) => ({ ...prev, state: 'PAUSED' }));
      showToast('Recording paused (SIGSTOP)', 'info');
    } catch (err: any) {
      showToast(`Pause error: ${err?.message || err}`, 'error');
    }
  };

  const handleResume = async () => {
    try {
      await Bridge.resumeRecording();
      setStatus((prev) => ({ ...prev, state: 'REC' }));
      showToast('Recording resumed (SIGCONT)', 'info');
    } catch (err: any) {
      showToast(`Resume error: ${err?.message || err}`, 'error');
    }
  };

  const handleToggle = async () => {
    if (status.state === 'REC') {
      await handlePause();
    } else if (status.state === 'PAUSED') {
      await handleResume();
    }
  };

  const handleStop = async () => {
    try {
      setStatus((prev) => ({ ...prev, state: 'PROCESSING' }));
      const savedPath = await Bridge.stopRecording();
      
      // Refresh recordings list
      const updatedList = await Bridge.getRecentRecordings();
      setRecentFiles(updatedList);
      
      setStatus({ state: 'IDLE', durationSecs: 0, formatted: '00:00:00' });
      showToast(`Saved lesson to MP4: ${savedPath}`, 'success');

      try {
        confetti({
          particleCount: 50,
          spread: 60,
          origin: { y: 0.85 },
          colors: ['#38BDF8', '#06B6D4', '#10B981', '#F43F5E'],
        });
      } catch (_) {}
    } catch (err: any) {
      showToast(`Stop error: ${err?.message || err}`, 'error');
      setStatus((prev) => ({ ...prev, state: 'ERROR', errorMessage: String(err) }));
    }
  };

  const handleSelectFolder = async () => {
    try {
      const selected = await Bridge.selectOutputDir();
      if (selected) {
        setConfig((prev) => ({ ...prev, outputDir: selected }));
      }
    } catch (err) {
      console.error('Select folder cancelled or error', err);
    }
  };

  const handleOpenFolder = async (path?: string) => {
    await Bridge.openFolder(path || config.outputDir);
  };

  const handlePlayVideo = async (path: string) => {
    await Bridge.playVideo(path);
  };

  // Keyboard Shortcuts (F9: Toggle Pause/Resume, F10: Stop Recording)
  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === 'F9') {
        e.preventDefault();
        handleToggle();
      } else if (e.key === 'F10') {
        e.preventDefault();
        if (status.state === 'REC' || status.state === 'PAUSED') {
          handleStop();
        }
      }
    },
    [status.state, config]
  );

  useEffect(() => {
    window.addEventListener('keydown', handleKeyDown);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [handleKeyDown]);

  const isBusy = status.state === 'REC' || status.state === 'PAUSED' || status.state === 'PROCESSING';

  return (
    <div className="flex flex-col h-screen w-screen bg-[#0B0F17] text-slate-100 select-none overflow-hidden">
      {/* Top Header */}
      <Header activeEncoder={activeEncoder} />

      {/* Main Content Area */}
      <main className="flex-1 overflow-y-auto p-6 space-y-5">
        {/* Status Display Card */}
        <StatusIndicator
          state={status.state}
          durationFormatted={status.formatted}
          errorMessage={status.errorMessage}
          resolution={config.resolution}
          fps={config.fps}
        />

        {/* Primary Controls */}
        <ControlPanel
          state={status.state}
          onStart={handleStart}
          onPause={handlePause}
          onResume={handleResume}
          onStop={handleStop}
        />

        {/* Video, Audio, and Output Settings */}
        <SettingsPanel
          config={config}
          onChange={(patch) => setConfig((prev) => ({ ...prev, ...patch }))}
          onSelectFolder={handleSelectFolder}
          onOpenFolder={() => handleOpenFolder()}
          disabled={isBusy}
          mics={audioDevices.mics}
          monitors={audioDevices.monitors}
        />

        {/* Recent Recordings List */}
        <RecentRecordings
          recordings={recentFiles}
          onPlay={handlePlayVideo}
          onOpenFolder={handleOpenFolder}
        />
      </main>

      {/* Footer Hotkey Bar & Notification Toast */}
      <footer className="relative p-4 pt-0">
        {notification && (
          <div className="absolute bottom-16 left-1/2 -translate-x-1/2 px-4 py-2 rounded-xl bg-slate-900/95 border border-slate-700 text-xs font-medium text-white shadow-2xl shadow-black/80 flex items-center gap-2 animate-fade-in z-50">
            <span className={`w-2 h-2 rounded-full ${
              notification.type === 'success' ? 'bg-emerald-400' : notification.type === 'error' ? 'bg-rose-400' : 'bg-cyan-400'
            }`} />
            <span>{notification.message}</span>
          </div>
        )}
        <HotkeyBar />
      </footer>
    </div>
  );
};

export default App;
