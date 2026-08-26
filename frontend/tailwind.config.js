/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: '#0B0F17',
        card: '#131B2A',
        'card-border': '#1E293B',
        accent: {
          blue: '#38BDF8',
          cyan: '#06B6D4',
          emerald: '#10B981',
          rose: '#F43F5E',
          amber: '#F59E0B',
          violet: '#8B5CF6'
        }
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
      },
      animation: {
        'pulse-fast': 'pulse 1.2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'glow-red': 'glowRed 2s ease-in-out infinite alternate',
        'glow-amber': 'glowAmber 2s ease-in-out infinite alternate',
        'wave': 'wave 1.5s ease-in-out infinite',
      },
      keyframes: {
        glowRed: {
          '0%': { boxShadow: '0 0 10px rgba(244, 63, 94, 0.4)' },
          '100%': { boxShadow: '0 0 25px rgba(244, 63, 94, 0.8), 0 0 40px rgba(244, 63, 94, 0.3)' },
        },
        glowAmber: {
          '0%': { boxShadow: '0 0 10px rgba(245, 158, 11, 0.4)' },
          '100%': { boxShadow: '0 0 25px rgba(245, 158, 11, 0.8), 0 0 40px rgba(245, 158, 11, 0.3)' },
        },
        wave: {
          '0%, 100%': { height: '8px' },
          '50%': { height: '28px' },
        }
      }
    },
  },
  plugins: [],
}
