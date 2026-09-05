import type { Block } from "@/lib/types";

type Props = { blocks: Block[] };

export function ArticleRenderer({ blocks }: Props) {
  return (
    <div className="prose max-w-none space-y-4">
      {blocks.map((block, i) => {
        if (block.type === "h") {
          const Tag = block.level <= 2 ? "h2" : "h3";
          return (
            <Tag
              key={i}
              className={
                block.level <= 2
                  ? "text-xl font-semibold text-fg"
                  : "text-lg font-medium text-fg"
              }
            >
              {block.text}
            </Tag>
          );
        }
        if (block.type === "code") {
          return (
            <pre
              key={i}
              className="overflow-x-auto rounded-relp-md border border-border bg-bg p-4 font-mono text-sm leading-relaxed text-fg"
            >
              <code>{block.text}</code>
            </pre>
          );
        }
        const isUrl = block.text.startsWith("http");
        if (isUrl) {
          return (
            <p key={i}>
              <a
                href={block.text}
                target="_blank"
                rel="noopener noreferrer"
                className="break-all text-primary underline-offset-2 hover:underline"
              >
                {block.text}
              </a>
            </p>
          );
        }
        return (
          <p key={i} className="leading-relaxed text-muted">
            {block.text}
          </p>
        );
      })}
    </div>
  );
}
