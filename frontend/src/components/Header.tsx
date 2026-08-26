import React from 'react';
import { Cpu, ShieldCheck, Zap } from 'lucide-react';
import { EncoderInfo } from '../types';

interface HeaderProps {
  activeEncoder?: EncoderInfo | null;
  isWayland?: boolean;
}

export const Header: React.FC<HeaderProps> = ({ activeEncoder }) => {
  return (
    <header className="flex items-center justify-between px-6 py-4 border-b border-slate-800/80 bg-slate-950/40 backdrop-blur-md">
      <div className="flex items-center gap-3">
        <div className="relative flex items-center justify-center w-11 h-11 rounded-2xl overflow-hidden shadow-lg shadow-cyan-500/20 border border-cyan-500/30">
          <img src="/logo.png" alt="GoRecorder Logo" className="w-full h-full object-cover" />
          <span className="absolute top-1 right-1 flex h-2.5 w-2.5">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-rose-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-rose-500"></span>
          </span>
        </div>
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-xl font-bold tracking-tight text-white flex items-center gap-1.5">
              Go<span className="text-transparent bg-clip-text bg-gradient-to-r from-cyan-400 to-blue-400">Recorder</span>
            </h1>
            <span className="px-2 py-0.5 text-[10px] font-semibold tracking-wider uppercase rounded-full bg-cyan-500/10 text-cyan-400 border border-cyan-500/20">
              v2.0 Native
            </span>
          </div>
          <p className="text-xs text-slate-400">Ubuntu 24.04 Wayland & PipeWire High-Performance Screen Recorder</p>
        </div>
      </div>

      <div className="flex items-center gap-3">
        {activeEncoder && (
          <div className="hidden sm:flex items-center gap-2 px-3 py-1.5 rounded-lg bg-slate-900/80 border border-slate-800 text-xs">
            {activeEncoder.isHardware ? (
              <Zap className="w-3.5 h-3.5 text-amber-400" />
            ) : (
              <Cpu className="w-3.5 h-3.5 text-blue-400" />
            )}
            <span className="text-slate-300 font-medium">{activeEncoder.name}</span>
          </div>
        )}

        <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-xs font-medium">
          <ShieldCheck className="w-3.5 h-3.5" />
          <span>Offline & Secure</span>
        </div>
      </div>
    </header>
  );
};
