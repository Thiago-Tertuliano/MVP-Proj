import type { Config } from "tailwindcss";

/** HSL channel + alpha do Tailwind (`bg-done/40`). Hex só em globals.css. */
const hsl = (name: string) => `hsl(var(${name}) / <alpha-value>)`;

const config: Config = {
  darkMode: ["class"],
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./stories/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        /* shadcn */
        border: hsl("--border"),
        input: hsl("--input"),
        ring: hsl("--ring"),
        background: hsl("--background"),
        foreground: hsl("--foreground"),
        primary: {
          DEFAULT: hsl("--primary"),
          foreground: hsl("--primary-foreground"),
          hover: hsl("--primary-hover"),
          muted: hsl("--primary-muted"),
        },
        secondary: {
          DEFAULT: hsl("--secondary"),
          foreground: hsl("--secondary-foreground"),
        },
        destructive: {
          DEFAULT: hsl("--destructive"),
          foreground: hsl("--destructive-foreground"),
        },
        muted: {
          DEFAULT: hsl("--muted"),
          foreground: hsl("--muted-foreground"),
        },
        accent: {
          DEFAULT: hsl("--accent"),
          foreground: hsl("--accent-foreground"),
          muted: hsl("--accent"),
        },
        popover: {
          DEFAULT: hsl("--popover"),
          foreground: hsl("--popover-foreground"),
        },
        card: {
          DEFAULT: hsl("--card"),
          foreground: hsl("--card-foreground"),
        },
        /* Relp — educação / status / feedback */
        learn: {
          DEFAULT: hsl("--relp-learn"),
          muted: hsl("--relp-learn-muted"),
        },
        done: {
          DEFAULT: hsl("--relp-done"),
          muted: hsl("--relp-done-muted"),
        },
        "progress-track": hsl("--relp-progress-track"),
        status: {
          draft: hsl("--relp-status-draft"),
          review: hsl("--relp-status-review"),
          published: hsl("--relp-status-published"),
          archived: hsl("--relp-status-archived"),
          locked: hsl("--relp-status-locked"),
        },
        danger: {
          DEFAULT: hsl("--relp-danger"),
          muted: hsl("--relp-danger-muted"),
        },
        warning: {
          DEFAULT: hsl("--relp-warning"),
          muted: hsl("--relp-warning-muted"),
        },
        info: {
          DEFAULT: hsl("--relp-info"),
          muted: hsl("--relp-info-muted"),
        },
        notify: {
          DEFAULT: hsl("--relp-notify"),
          unread: hsl("--relp-notify-unread"),
        },
        /* Aliases de superfície (mock / páginas) */
        bg: hsl("--background"),
        fg: hsl("--foreground"),
        surface: {
          DEFAULT: hsl("--card"),
          raised: hsl("--card"),
          border: hsl("--border"),
        },
        success: hsl("--relp-done"),
      },
      spacing: {
        "relp-1": "var(--relp-space-1)",
        "relp-2": "var(--relp-space-2)",
        "relp-3": "var(--relp-space-3)",
        "relp-4": "var(--relp-space-4)",
        "relp-6": "var(--relp-space-6)",
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
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
  plugins: [require("tailwindcss-animate")],
};

export default config;
