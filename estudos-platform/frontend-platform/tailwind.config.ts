import type { Config } from "tailwindcss";

/** Cor tokenizada: rgb + alpha do Tailwind (`bg-primary/40`). Hex só em globals.css. */
const token = (name: string) => `rgb(var(${name}) / <alpha-value>)`;

const config: Config = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: token("--relp-primary"),
          hover: token("--relp-primary-hover"),
          muted: token("--relp-primary-muted"),
          fg: token("--relp-primary-fg"),
        },
        learn: {
          DEFAULT: token("--relp-learn"),
          muted: token("--relp-learn-muted"),
        },
        done: {
          DEFAULT: token("--relp-done"),
          muted: token("--relp-done-muted"),
        },
        "progress-track": token("--relp-progress-track"),
        status: {
          draft: token("--relp-status-draft"),
          review: token("--relp-status-review"),
          published: token("--relp-status-published"),
          archived: token("--relp-status-archived"),
          locked: token("--relp-status-locked"),
        },
        danger: {
          DEFAULT: token("--relp-danger"),
          muted: token("--relp-danger-muted"),
        },
        warning: {
          DEFAULT: token("--relp-warning"),
          muted: token("--relp-warning-muted"),
        },
        info: {
          DEFAULT: token("--relp-info"),
          muted: token("--relp-info-muted"),
        },
        notify: {
          DEFAULT: token("--relp-notify"),
          unread: token("--relp-notify-unread"),
        },
        bg: token("--relp-bg"),
        fg: token("--relp-fg"),
        muted: token("--relp-muted"),
        surface: {
          DEFAULT: token("--relp-surface"),
          raised: token("--relp-surface-raised"),
        },
        border: token("--relp-border"),
        /* Aliases legados → tokens Relp (remover após F2) */
        accent: {
          DEFAULT: token("--relp-primary"),
          muted: token("--relp-primary-muted"),
        },
        success: token("--relp-done"),
      },
      spacing: {
        "relp-1": "var(--relp-space-1)",
        "relp-2": "var(--relp-space-2)",
        "relp-3": "var(--relp-space-3)",
        "relp-4": "var(--relp-space-4)",
        "relp-6": "var(--relp-space-6)",
      },
      borderRadius: {
        "relp-sm": "var(--relp-radius-sm)",
        "relp-md": "var(--relp-radius-md)",
        "relp-lg": "var(--relp-radius-lg)",
      },
      fontFamily: {
        sans: ["var(--font-geist-sans)", "system-ui", "sans-serif"],
        mono: ["var(--font-geist-mono)", "monospace"],
      },
    },
  },
  plugins: [],
};

export default config;
