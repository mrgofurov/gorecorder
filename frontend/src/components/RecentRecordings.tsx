import React from 'react';
import { Play, FolderOpen, Film, Clock, HardDrive } from 'lucide-react';
import { RecordedFileEntry } from '../types';

interface RecentRecordingsProps {
  recordings: RecordedFileEntry[];
  onPlay: (path: string) => void;
  onOpenFolder: (path: string) => void;
}

export const RecentRecordings: React.FC<RecentRecordingsProps> = ({
  recordings,
  onPlay,
  onOpenFolder,
}) => {
  if (recordings.length === 0) {
    return (
      <div className="p-4 rounded-xl glass-panel-subtle text-center text-xs text-slate-500 py-6">
        No recordings yet in this session. Start recording to produce lessons!
      </div>
    );
  }

  const formatDuration = (secs: number) => {
    const m = Math.floor(secs / 60);
    const s = secs % 60;
    return `${m}:${s.toString().padStart(2, '0')}`;
  };

  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center justify-between px-1">
        <h3 className="text-xs font-semibold text-slate-400 uppercase tracking-wider flex items-center gap-1.5">
          <Film className="w-3.5 h-3.5 text-cyan-400" />
          <span>Recent Lessons</span>
        </h3>
        <span className="text-[11px] text-slate-500">{recordings.length} items</span>
      </div>

      <div className="space-y-2 max-h-36 overflow-y-auto pr-1">
        {recordings.map((rec, idx) => (
          <div
            key={idx}
            className="flex items-center justify-between p-2.5 rounded-lg bg-slate-900/60 hover:bg-slate-900/90 border border-slate-800/80 transition-colors group"
          >
            <div className="flex items-center gap-2.5 min-w-0">
              <div className="p-1.5 rounded bg-cyan-500/10 text-cyan-400">
                <Film className="w-3.5 h-3.5" />
              </div>
              <div className="min-w-0">
                <div className="text-xs font-medium text-slate-200 truncate">{rec.filename}</div>
                <div className="flex items-center gap-3 text-[10px] text-slate-400 mt-0.5">
                  <span className="flex items-center gap-1">
                    <Clock className="w-3 h-3 text-slate-500" />
                    {formatDuration(rec.durationSec)}
                  </span>
                  {rec.sizeHuman && (
                    <span className="flex items-center gap-1">
                      <HardDrive className="w-3 h-3 text-slate-500" />
                      {rec.sizeHuman}
                    </span>
                  )}
                </div>
              </div>
            </div>

            <div className="flex items-center gap-1.5 opacity-80 group-hover:opacity-100 transition-opacity">
              <button
                type="button"
                onClick={() => onPlay(rec.filePath)}
                className="p-1.5 rounded-md hover:bg-cyan-500/20 text-slate-300 hover:text-cyan-400 transition-colors cursor-pointer"
                title="Play Video"
              >
                <Play className="w-3.5 h-3.5 fill-current" />
              </button>
              <button
                type="button"
                onClick={() => onOpenFolder(rec.filePath)}
                className="p-1.5 rounded-md hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors cursor-pointer"
                title="Show in Folder"
              >
                <FolderOpen className="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
