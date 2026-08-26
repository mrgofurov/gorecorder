import React from 'react';
import { Radio, Pause, CheckCircle2, AlertCircle, Loader2, Volume2 } from 'lucide-react';
import { RecordingState } from '../types';

interface StatusIndicatorProps {
  state: RecordingState;
  durationFormatted: string;
  errorMessage?: string;
  resolution: string;
  fps: number;
}

export const StatusIndicator: React.FC<StatusIndicatorProps> = ({
  state,
  durationFormatted,
  errorMessage,
  resolution,
  fps,
}) => {
  const isRecording = state === 'REC';
  const isPaused = state === 'PAUSED';
  const isProcessing = state === 'PROCESSING';

  return (
    <div className="relative overflow-hidden rounded-2xl glass-panel p-6">
      {/* Background ambient glow according to state */}
      {isRecording && (
        <div className="absolute -top-24 -left-24 w-64 h-64 bg-rose-500/15 rounded-full blur-3xl pointer-events-none animate-pulse" />
      )}
      {isPaused && (
        <div className="absolute -top-24 -left-24 w-64 h-64 bg-amber-500/15 rounded-full blur-3xl pointer-events-none" />
      )}

      <div className="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
        {/* State Badge & Information */}
        <div className="flex items-center gap-4">
          <div className="relative">
            {state === 'REC' && (
              <div className="flex items-center gap-2 px-4 py-2 rounded-xl bg-rose-500/15 border border-rose-500/30 text-rose-400 font-bold tracking-wider shadow-lg shadow-rose-500/10">
                <span className="relative flex h-3 w-3">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-rose-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-3 w-3 bg-rose-500"></span>
                </span>
                <Radio className="w-4 h-4 animate-pulse" />
                <span>REC</span>
              </div>
            )}

            {state === 'PAUSED' && (
              <div className="flex items-center gap-2 px-4 py-2 rounded-xl bg-amber-500/15 border border-amber-500/30 text-amber-400 font-bold tracking-wider shadow-lg shadow-amber-500/10">
                <Pause className="w-4 h-4 animate-bounce" />
                <span>PAUSED</span>
              </div>
            )}

            {state === 'PROCESSING' && (
              <div className="flex items-center gap-2 px-4 py-2 rounded-xl bg-blue-500/15 border border-blue-500/30 text-blue-400 font-bold tracking-wider">
                <Loader2 className="w-4 h-4 animate-spin" />
                <span>REMUXING MP4</span>
              </div>
            )}

            {state === 'IDLE' && (
              <div className="flex items-center gap-2 px-4 py-2 rounded-xl bg-slate-800/80 border border-slate-700 text-slate-300 font-medium">
                <CheckCircle2 className="w-4 h-4 text-emerald-400" />
                <span>READY</span>
              </div>
            )}

            {state === 'ERROR' && (
              <div className="flex items-center gap-2 px-4 py-2 rounded-xl bg-red-950/60 border border-red-800/60 text-red-400 font-medium">
                <AlertCircle className="w-4 h-4" />
                <span>ERROR</span>
              </div>
            )}
          </div>

          <div>
            <div className="text-xs text-slate-400 uppercase tracking-wider font-semibold">
              {isRecording ? 'Live Session' : isPaused ? 'Session Paused' : isProcessing ? 'Finalizing MP4' : 'Capture Mode'}
            </div>
            <div className="text-sm font-medium text-slate-200 flex items-center gap-2 mt-0.5">
              <span>{resolution.toUpperCase()}</span>
              <span className="text-slate-600">•</span>
              <span>{fps} FPS</span>
            </div>
          </div>
        </div>

        {/* Digital Clock & Audio visualizer */}
        <div className="flex items-center gap-6">
          {(isRecording || isPaused) && (
            <div className="hidden sm:flex items-center gap-1.5 h-8 px-3 rounded-lg bg-slate-900/60 border border-slate-800">
              <Volume2 className={`w-4 h-4 ${isRecording ? 'text-cyan-400 animate-pulse' : 'text-slate-500'}`} />
              <div className="flex items-end gap-1 h-5 w-16 px-1">
                {[...Array(6)].map((_, i) => (
                  <span
                    key={i}
                    className={`w-1 rounded-full transition-all duration-300 ${
                      isRecording
                        ? 'bg-cyan-400 animate-pulse'
                        : 'bg-slate-700 h-1.5'
                    }`}
                    style={{
                      height: isRecording ? `${Math.max(6, (i + 1) * 4 * (i % 2 === 0 ? 0.9 : 1.2))}px` : '4px',
                      animationDelay: `${i * 150}ms`,
                    }}
                  />
                ))}
              </div>
            </div>
          )}

          <div className="flex flex-col items-end">
            <div className="text-[10px] text-slate-400 font-mono tracking-wider uppercase">Recording Time</div>
            <div className={`text-4xl sm:text-5xl font-mono font-bold tracking-tight ${
              isRecording ? 'text-white drop-shadow-[0_0_15px_rgba(255,255,255,0.3)]' : isPaused ? 'text-amber-400/90' : 'text-slate-400'
            }`}>
              {durationFormatted || '00:00:00'}
            </div>
          </div>
        </div>
      </div>

      {errorMessage && (
        <div className="mt-4 p-3 rounded-lg bg-red-950/40 border border-red-800/40 text-red-300 text-xs flex items-center gap-2">
          <AlertCircle className="w-4 h-4 flex-shrink-0" />
          <span>{errorMessage}</span>
        </div>
      )}
    </div>
  );
};
