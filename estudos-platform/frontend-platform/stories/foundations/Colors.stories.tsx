import type { Meta, StoryObj } from "@storybook/react";

const SWATCHES = [
  { name: "primary", className: "bg-primary" },
  { name: "secondary", className: "bg-secondary border border-border" },
  { name: "accent", className: "bg-accent" },
  { name: "destructive", className: "bg-destructive" },
  { name: "muted", className: "bg-muted border border-border" },
  { name: "background", className: "bg-background border border-border" },
  { name: "card", className: "bg-card border border-border" },
] as const;

function ColorSwatches() {
  return (
    <div className="grid max-w-xl gap-4 sm:grid-cols-2">
      {SWATCHES.map((swatch) => (
        <div key={swatch.name} className="space-y-2">
          <div className={`h-16 rounded-lg ${swatch.className}`} />
          <p className="text-sm text-foreground">{swatch.name}</p>
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
