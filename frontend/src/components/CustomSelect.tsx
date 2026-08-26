import React, { useState, useRef, useEffect } from 'react';
import { ChevronDown, Check } from 'lucide-react';

export interface Option<T> {
  value: T;
  label: string;
  sublabel?: string;
}

interface CustomSelectProps<T> {
  value: T;
  options: Option<T>[];
  onChange: (value: T) => void;
  disabled?: boolean;
  icon?: React.ReactNode;
}

export function CustomSelect<T extends string | number>({
  value,
  options,
  onChange,
  disabled = false,
  icon,
}: CustomSelectProps<T>) {
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  const selectedOption = options.find((opt) => opt.value === value) || options[0];

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  return (
    <div className={`relative w-full ${isOpen ? 'z-50' : 'z-10'}`} ref={containerRef}>
      <button
        type="button"
        disabled={disabled}
        onClick={() => !disabled && setIsOpen(!isOpen)}
        className={`w-full flex items-center justify-between px-3.5 py-2.5 rounded-xl bg-slate-900/90 border transition-all text-left cursor-pointer ${
          isOpen
            ? 'border-cyan-500 shadow-lg shadow-cyan-500/10 ring-1 ring-cyan-500'
            : 'border-slate-700/80 hover:border-slate-600 hover:bg-slate-900'
        } disabled:opacity-50 disabled:cursor-not-allowed`}
      >
        <div className="flex items-center gap-2 min-w-0">
          {icon && <span className="text-cyan-400 flex-shrink-0">{icon}</span>}
          <div className="min-w-0 truncate">
            <span className="text-sm font-medium text-slate-100">{selectedOption?.label}</span>
            {selectedOption?.sublabel && (
              <span className="ml-1.5 text-xs text-slate-400 hidden sm:inline">
                ({selectedOption.sublabel})
              </span>
            )}
          </div>
        </div>
        <ChevronDown
          className={`w-4 h-4 text-slate-400 transition-transform duration-200 flex-shrink-0 ${
            isOpen ? 'rotate-180 text-cyan-400' : ''
          }`}
        />
      </button>

      {isOpen && (
        <div className="absolute z-50 left-0 right-0 top-full mt-1.5 py-1.5 rounded-xl bg-[#131B2A] border border-slate-700 shadow-2xl shadow-black max-h-60 overflow-y-auto animate-in fade-in slide-in-from-top-2 duration-150">
          {options.map((opt) => {
            const isSelected = opt.value === value;
            return (
              <button
                key={String(opt.value)}
                type="button"
                onClick={() => {
                  onChange(opt.value);
                  setIsOpen(false);
                }}
                className={`w-full flex items-center justify-between px-3.5 py-2.5 text-left text-sm transition-colors cursor-pointer ${
                  isSelected
                    ? 'bg-cyan-500/20 text-cyan-300 font-semibold'
                    : 'text-slate-200 hover:bg-slate-800 hover:text-white'
                }`}
              >
                <div className="flex flex-col min-w-0">
                  <span className="truncate">{opt.label}</span>
                  {opt.sublabel && (
                    <span className="text-[11px] text-slate-400 font-normal truncate">
                      {opt.sublabel}
                    </span>
                  )}
                </div>
                {isSelected && <Check className="w-4 h-4 text-cyan-400 ml-2 flex-shrink-0" />}
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}
