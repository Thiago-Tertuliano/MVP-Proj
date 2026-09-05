import React from "react";
import type { Preview } from "@storybook/react";
import { withThemeByClassName } from "@storybook/addon-themes";

import "../app/globals.css";

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
    layout: "padded",
    a11y: {
      test: "todo",
    },
    backgrounds: {
      disable: true,
    },
  },
  decorators: [
    withThemeByClassName({
      themes: {
        light: "",
        dark: "dark",
      },
      defaultTheme: "light",
    }),
    (Story) =>
      React.createElement(
        "div",
        {
          className:
            "min-h-screen bg-surface p-6 font-sans text-[var(--foreground)] antialiased",
        },
        React.createElement(Story),
      ),
  ],
};

export default preview;
