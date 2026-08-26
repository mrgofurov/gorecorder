export type RecordingState = 'IDLE' | 'REC' | 'PAUSED' | 'PROCESSING' | 'ERROR';

export type ResolutionOption = 'native' | '480p' | '720p' | '1080p' | '1440p';

export type AudioSourceOption = 'none' | 'mic' | 'system' | 'both';

export interface RecordConfig {
  resolution: ResolutionOption;
  fps: number;
  audioSource: AudioSourceOption;
  outputDir: string;
  micDevice?: string;
  sysDevice?: string;
}

export interface EncoderInfo {
  type: string;
  name: string;
  description: string;
  devicePath?: string;
  isHardware: boolean;
}

export interface AudioDeviceInfo {
  id: string;
  name: string;
  description: string;
  isDefault: boolean;
  isMonitor: boolean;
}

export interface RecordingStatus {
  state: RecordingState;
  durationSecs: number;
  formatted: string;
  outputFile?: string;
  tempFile?: string;
  encoderUsed?: string;
  errorMessage?: string;
}

export interface RecordedFileEntry {
  filename: string;
  filePath: string;
  sizeHuman: string;
  durationSec: number;
  createdAt: string;
}
