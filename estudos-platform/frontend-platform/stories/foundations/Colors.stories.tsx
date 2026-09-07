import type { Meta, StoryObj } from "@storybook/react";

const SWATCHES = [
  { name: "accent", className: "bg-accent" },
  { name: "accent-muted", className: "bg-accent-muted" },
  { name: "success", className: "bg-success" },
  { name: "warning", className: "bg-warning" },
  { name: "surface", className: "bg-surface border border-surface-border" },
  { name: "surface-raised", className: "bg-surface-raised border border-surface-border" },
] as const;

function ColorSwatches() {
  return (
    <div className="grid max-w-xl gap-4 sm:grid-cols-2">
      {SWATCHES.map((swatch) => (
        <div key={swatch.name} className="space-y-2">
          <div className={`h-16 rounded-lg ${swatch.className}`} />
          <p className="text-sm text-[var(--foreground)]">{swatch.name}</p>
        </div>
      ))}
    </div>
  );
}

const meta = {
  title: "Foundations/Colors",
  component: ColorSwatches,
  parameters: {
    layout: "padded",
  },
} satisfies Meta<typeof ColorSwatches>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
