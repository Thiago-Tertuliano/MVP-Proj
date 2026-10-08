import type { StorybookConfig } from "@storybook/nextjs";

const config: StorybookConfig = {
  stories: [
    "../stories/**/*.mdx",
    "../stories/**/*.stories.@(js|jsx|mjs|ts|tsx)",
    "../components/**/*.stories.@(js|jsx|mjs|ts|tsx)",
  ],
  addons: [
    "@storybook/addon-essentials",
    "@storybook/addon-a11y",
    "@storybook/addon-themes",
    "@storybook/addon-interactions",
    // Estados :hover / :focus / :active em stories (parameters.pseudo).
    "storybook-addon-pseudo-states",
  ],
  framework: {
    name: "@storybook/nextjs",
    options: {},
  },
};

export default config;
