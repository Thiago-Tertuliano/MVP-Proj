import type { Block } from "@/lib/types";

type Props = { blocks: Block[] };

export function ArticleRenderer({ blocks }: Props) {
  return (
    <div className="prose prose-invert max-w-none space-y-4">
      {blocks.map((block, i) => {
        if (block.type === "h") {
          const Tag = block.level <= 2 ? "h2" : "h3";
          return (
            <Tag
              key={i}
              className={
                block.level <= 2
                  ? "text-xl font-semibold text-white"
                  : "text-lg font-medium text-zinc-200"
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
              className="overflow-x-auto rounded-lg border border-surface-border bg-[#0d1117] p-4 font-mono text-sm leading-relaxed text-zinc-200"
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
                className="break-all text-accent underline-offset-2 hover:underline"
              >
                {block.text}
              </a>
            </p>
          );
        }
        return (
          <p key={i} className="leading-relaxed text-zinc-300">
            {block.text}
          </p>
        );
      })}
    </div>
  );
}
