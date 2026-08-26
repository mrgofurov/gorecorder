import React from 'react';
import { Monitor, Gauge, Mic, Folder, FolderOpen } from 'lucide-react';
import { AudioDeviceInfo, AudioSourceOption, RecordConfig, ResolutionOption } from '../types';
import { CustomSelect, Option } from './CustomSelect';

interface SettingsPanelProps {
  config: RecordConfig;
  onChange: (cfg: Partial<RecordConfig>) => void;
  onSelectFolder: () => void;
  onOpenFolder: () => void;
  disabled?: boolean;
  mics?: AudioDeviceInfo[];
  monitors?: AudioDeviceInfo[];
}

const resolutionOptions: Option<ResolutionOption>[] = [
  { value: 'native', label: 'Native Display', sublabel: 'Original' },
  { value: '1080p', label: '1080p Full HD', sublabel: '1920x1080' },
  { value: '720p', label: '720p HD', sublabel: '1280x720' },
  { value: '1440p', label: '1440p 2K', sublabel: '2560x1440' },
  { value: '480p', label: '480p SD', sublabel: '854x480' },
];

const fpsOptions: Option<number>[] = [
  { value: 60, label: '60 FPS', sublabel: 'Smooth & Crisp' },
  { value: 30, label: '30 FPS', sublabel: 'Standard' },
];

export const SettingsPanel: React.FC<SettingsPanelProps> = ({
  config,
  onChange,
  onSelectFolder,
  onOpenFolder,
  disabled = false,
}) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {/* Resolution & FPS */}
      <div className="relative z-30 p-4 rounded-xl glass-panel flex flex-col gap-4">
        <div className="flex items-center gap-2 text-xs font-semibold text-cyan-400 uppercase tracking-wider">
          <Monitor className="w-4 h-4" />
          <span>Video Configuration</span>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label className="block text-xs text-slate-400 mb-1.5 font-medium">Resolution</label>
            <CustomSelect
              value={config.resolution}
              options={resolutionOptions}
              onChange={(val) => onChange({ resolution: val })}
              disabled={disabled}
            />
          </div>

          <div>
            <label className="block text-xs text-slate-400 mb-1.5 font-medium flex items-center gap-1">
              <Gauge className="w-3 h-3 text-cyan-400" />
              <span>Framerate</span>
            </label>
            <CustomSelect
              value={config.fps}
              options={fpsOptions}
              onChange={(val) => onChange({ fps: val })}
              disabled={disabled}
            />
          </div>
        </div>
      </div>

      {/* Audio Source Selector */}
      <div className="relative z-20 p-4 rounded-xl glass-panel flex flex-col gap-4">
        <div className="flex items-center gap-2 text-xs font-semibold text-cyan-400 uppercase tracking-wider">
          <Mic className="w-4 h-4" />
          <span>Audio Sources</span>
        </div>

        <div>
          <label className="block text-xs text-slate-400 mb-1.5 font-medium">Audio Capture Mode (AAC 192k)</label>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-2">
            {[
              { id: 'both', label: 'Mic + System', desc: 'Combined' },
              { id: 'mic', label: 'Microphone', desc: 'Voice Only' },
              { id: 'system', label: 'System Audio', desc: 'Desktop sounds' },
              { id: 'none', label: 'Muted', desc: 'No Audio' },
            ].map((opt) => (
              <button
                key={opt.id}
                type="button"
                disabled={disabled}
                onClick={() => onChange({ audioSource: opt.id as AudioSourceOption })}
                className={`flex flex-col items-center justify-center p-2.5 rounded-xl border text-center transition-all cursor-pointer ${
                  config.audioSource === opt.id
                    ? 'bg-cyan-500/15 border-cyan-500/50 text-cyan-300 font-semibold shadow-md shadow-cyan-500/10 ring-1 ring-cyan-500/50'
                    : 'bg-slate-900/60 border-slate-800 text-slate-400 hover:text-slate-200 hover:border-slate-700 hover:bg-slate-900'
                } disabled:opacity-50`}
              >
                <span className="text-xs">{opt.label}</span>
                <span className="text-[10px] text-slate-400 font-normal mt-0.5">{opt.desc}</span>
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Destination Folder Selector */}
      <div className="relative z-10 md:col-span-2 p-4 rounded-xl glass-panel flex flex-col gap-2">
        <div className="flex items-center justify-between">
          <label className="text-xs font-semibold text-cyan-400 uppercase tracking-wider flex items-center gap-2">
            <Folder className="w-4 h-4" />
            <span>Output Directory (MP4)</span>
          </label>

          <button
            type="button"
            onClick={onOpenFolder}
            className="flex items-center gap-1.5 text-xs text-slate-400 hover:text-cyan-400 transition-colors cursor-pointer"
            title="Open folder in File Manager"
          >
            <FolderOpen className="w-3.5 h-3.5" />
            <span>Open Folder</span>
          </button>
        </div>

        <div className="flex items-center gap-2">
          <div className="flex-1 bg-slate-900/90 border border-slate-700/80 rounded-xl px-3.5 py-2.5 text-xs font-mono text-slate-300 truncate">
            {config.outputDir || '~/Videos'}
          </div>

          <button
            type="button"
            onClick={onSelectFolder}
            disabled={disabled}
            className="px-4 py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 border border-slate-700 text-xs font-medium text-slate-200 transition-colors disabled:opacity-50 cursor-pointer"
          >
            Browse...
          </button>
        </div>
      </div>
    </div>
  );
};
