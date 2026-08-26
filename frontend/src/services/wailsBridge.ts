import { AudioDeviceInfo, EncoderInfo, RecordConfig, RecordedFileEntry, RecordingStatus } from '../types';
import * as AppMethods from '../../wailsjs/go/main/App';
import * as Runtime from '../../wailsjs/runtime/runtime';

export const Bridge = {
  async startRecording(cfg: RecordConfig): Promise<void> {
    try {
      if (AppMethods.StartRecording) {
        return await AppMethods.StartRecording(cfg as any);
      }
    } catch (e) {
      console.warn('Wails IPC StartRecording fallback', e);
    }
  },

  async pauseRecording(): Promise<void> {
    try {
      if (AppMethods.PauseRecording) {
        return await AppMethods.PauseRecording();
      }
    } catch (e) {
      console.warn('Wails IPC PauseRecording fallback', e);
    }
  },

  async resumeRecording(): Promise<void> {
    try {
      if (AppMethods.ResumeRecording) {
        return await AppMethods.ResumeRecording();
      }
    } catch (e) {
      console.warn('Wails IPC ResumeRecording fallback', e);
    }
  },

  async togglePauseResume(): Promise<string> {
    try {
      if (AppMethods.TogglePauseResume) {
        return await AppMethods.TogglePauseResume();
      }
    } catch (e) {
      console.warn('Wails IPC TogglePauseResume fallback', e);
    }
    return 'REC';
  },

  async stopRecording(): Promise<string> {
    try {
      if (AppMethods.StopRecording) {
        return await AppMethods.StopRecording();
      }
    } catch (e) {
      console.warn('Wails IPC StopRecording fallback', e);
    }
    return 'golang-lesson.mp4';
  },

  async getStatus(): Promise<RecordingStatus> {
    try {
      if (AppMethods.GetStatus) {
        const res = await AppMethods.GetStatus();
        if (res) return res as RecordingStatus;
      }
    } catch (e) {
      console.warn('Wails IPC GetStatus fallback', e);
    }
    return {
      state: 'IDLE',
      durationSecs: 0,
      formatted: '00:00:00',
      encoderUsed: 'VAAPI Hardware (Intel iGPU)',
    };
  },

  async getAudioDevices(): Promise<{ mics: AudioDeviceInfo[]; monitors: AudioDeviceInfo[] }> {
    try {
      if (AppMethods.GetAudioDevices) {
        const res = await AppMethods.GetAudioDevices();
        if (res && res.mics) return res;
      }
    } catch (e) {
      console.warn('Wails IPC GetAudioDevices fallback', e);
    }
    return {
      mics: [{ id: 'default', name: 'default', description: 'Default Microphone', isDefault: true, isMonitor: false }],
      monitors: [{ id: 'default.monitor', name: 'default.monitor', description: 'Default System Audio', isDefault: true, isMonitor: true }],
    };
  },

  async getEncoders(): Promise<EncoderInfo[]> {
    try {
      if (AppMethods.GetEncoders) {
        const res = await AppMethods.GetEncoders();
        if (res && res.length > 0) return res as EncoderInfo[];
      }
    } catch (e) {
      console.warn('Wails IPC GetEncoders fallback', e);
    }
    return [
      { type: 'h264_vaapi', name: 'Intel/AMD VAAPI (Hardware)', description: 'Hardware acceleration via VA-API', isHardware: true },
      { type: 'libx264', name: 'libx264 (Software CPU)', description: 'Fast CPU H.264 encoder', isHardware: false },
    ];
  },

  async getRecentRecordings(): Promise<RecordedFileEntry[]> {
    try {
      if (AppMethods.GetRecentRecordings) {
        const res = await AppMethods.GetRecentRecordings();
        if (res) return res as RecordedFileEntry[];
      }
    } catch (e) {
      console.warn('Wails IPC GetRecentRecordings fallback', e);
    }
    return [];
  },

  async getDefaultOutputDir(): Promise<string> {
    try {
      if (AppMethods.GetDefaultOutputDir) {
        const res = await AppMethods.GetDefaultOutputDir();
        if (res) return res;
      }
    } catch (e) {
      console.warn('Wails IPC GetDefaultOutputDir fallback', e);
    }
    return '~/Videos';
  },

  async selectOutputDir(): Promise<string> {
    try {
      if (AppMethods.SelectOutputDir) {
        const res = await AppMethods.SelectOutputDir();
        if (res) return res;
      }
    } catch (e) {
      console.warn('Wails IPC SelectOutputDir fallback', e);
    }
    return '~/Videos';
  },

  async openFolder(path: string): Promise<void> {
    try {
      if (AppMethods.OpenFolder) {
        return await AppMethods.OpenFolder(path);
      }
    } catch (e) {
      console.warn('Wails IPC OpenFolder fallback', e);
    }
  },

  async playVideo(path: string): Promise<void> {
    try {
      if (AppMethods.PlayVideo) {
        return await AppMethods.PlayVideo(path);
      }
    } catch (e) {
      console.warn('Wails IPC PlayVideo fallback', e);
    }
  },

  onStatusChange(callback: (status: RecordingStatus) => void): () => void {
    try {
      if (Runtime && Runtime.EventsOn) {
        return Runtime.EventsOn('recorder:status', callback);
      }
    } catch (e) {
      console.warn('Wails runtime EventsOn fallback', e);
    }
    return () => {};
  },
};
