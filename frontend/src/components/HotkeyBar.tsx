import React from 'react';
import { Keyboard } from 'lucide-react';

export const HotkeyBar: React.FC = () => {
  return (
    <div className="flex items-center justify-between px-4 py-2 rounded-xl bg-slate-950/60 border border-slate-800/80 text-[11px] text-slate-400">
      <div className="flex items-center gap-2">
        <Keyboard className="w-3.5 h-3.5 text-cyan-400" />
        <span className="font-semibold text-slate-300">Global Hotkeys:</span>
      </div>

      <div className="flex items-center gap-4">
        <div className="flex items-center gap-1.5">
          <kbd className="px-2 py-0.5 rounded bg-slate-800 border border-slate-700 text-slate-200 font-mono font-bold text-[10px] shadow-sm">
            F9
          </kbd>
          <span>Pause / Resume</span>
        </div>

        <div className="flex items-center gap-1.5">
          <kbd className="px-2 py-0.5 rounded bg-slate-800 border border-slate-700 text-slate-200 font-mono font-bold text-[10px] shadow-sm">
            F10
          </kbd>
          <span>Stop & Save</span>
        </div>
      </div>
    </div>
  );
};
