import React from 'react';
import { Play, Pause, Square, Loader2, Sparkles } from 'lucide-react';
import { RecordingState } from '../types';

interface ControlPanelProps {
  state: RecordingState;
  onStart: () => void;
  onPause: () => void;
  onResume: () => void;
  onStop: () => void;
}

export const ControlPanel: React.FC<ControlPanelProps> = ({
  state,
  onStart,
  onPause,
  onResume,
  onStop,
}) => {
  const isIdle = state === 'IDLE' || state === 'ERROR';
  const isRecording = state === 'REC';
  const isPaused = state === 'PAUSED';
  const isProcessing = state === 'PROCESSING';

  return (
    <div className="flex items-center justify-center gap-4 py-2">
      {isIdle && (
        <button
          onClick={onStart}
          className="group relative flex items-center justify-center gap-3 px-10 py-4 rounded-2xl bg-gradient-to-r from-rose-600 to-red-500 hover:from-rose-500 hover:to-red-400 text-white font-bold text-base shadow-xl shadow-rose-600/30 hover:shadow-rose-600/50 hover:scale-[1.02] active:scale-[0.98] transition-all duration-200 cursor-pointer"
        >
          <span className="relative flex h-3.5 w-3.5">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-white opacity-75"></span>
            <span className="relative inline-flex rounded-full h-3.5 w-3.5 bg-white"></span>
          </span>
          <span>Start Recording</span>
          <Sparkles className="w-4 h-4 text-rose-200 opacity-80 group-hover:rotate-12 transition-transform" />
        </button>
      )}

      {(isRecording || isPaused) && (
        <>
          {/* Pause / Resume Button */}
          {isRecording ? (
            <button
              onClick={onPause}
              className="flex items-center gap-2.5 px-6 py-3.5 rounded-xl bg-amber-500/15 hover:bg-amber-500/25 border border-amber-500/30 text-amber-400 font-semibold text-sm shadow-lg shadow-amber-500/10 hover:scale-[1.02] active:scale-[0.98] transition-all duration-200 cursor-pointer"
            >
              <Pause className="w-4 h-4 fill-amber-400" />
              <span>Pause</span>
              <kbd className="ml-1.5 px-1.5 py-0.5 text-[10px] font-mono bg-amber-500/20 rounded border border-amber-500/30 text-amber-300">
                F9
              </kbd>
            </button>
          ) : (
            <button
              onClick={onResume}
              className="flex items-center gap-2.5 px-6 py-3.5 rounded-xl bg-emerald-500/15 hover:bg-emerald-500/25 border border-emerald-500/30 text-emerald-400 font-semibold text-sm shadow-lg shadow-emerald-500/10 hover:scale-[1.02] active:scale-[0.98] transition-all duration-200 cursor-pointer"
            >
              <Play className="w-4 h-4 fill-emerald-400" />
              <span>Resume</span>
              <kbd className="ml-1.5 px-1.5 py-0.5 text-[10px] font-mono bg-emerald-500/20 rounded border border-emerald-500/30 text-emerald-300">
                F9
              </kbd>
            </button>
          )}

          {/* Stop Recording Button */}
          <button
            onClick={onStop}
            className="flex items-center gap-2.5 px-7 py-3.5 rounded-xl bg-gradient-to-r from-red-600 to-rose-600 hover:from-red-500 hover:to-rose-500 text-white font-semibold text-sm shadow-lg shadow-red-600/30 hover:scale-[1.02] active:scale-[0.98] transition-all duration-200 cursor-pointer"
          >
            <Square className="w-4 h-4 fill-white" />
            <span>Stop & Save</span>
            <kbd className="ml-1.5 px-1.5 py-0.5 text-[10px] font-mono bg-black/30 rounded border border-white/20 text-white/90">
              F10
            </kbd>
          </button>
        </>
      )}

      {isProcessing && (
        <div className="flex items-center gap-3 px-8 py-3.5 rounded-xl bg-slate-800/80 border border-slate-700 text-slate-300 font-medium text-sm">
          <Loader2 className="w-4 h-4 animate-spin text-cyan-400" />
          <span>Processing and remuxing MP4...</span>
        </div>
      )}
    </div>
  );
};
