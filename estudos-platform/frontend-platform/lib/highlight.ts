/**
 * Realce de sintaxe mínimo, sem dependência: comentários, strings, números e palavras-chave
 * das linguagens que aparecem nas aulas (Go, JS/TS, Python, SQL, Bash). Linguagem desconhecida
 * devolve só texto simples — o bloco continua legível.
 */

export type TokenType = "plain" | "keyword" | "string" | "comment" | "number";
export type Token = { type: TokenType; text: string };

const KEYWORDS: Record<string, string[]> = {
  go: [
    "break", "case", "chan", "const", "continue", "default", "defer", "else", "fallthrough", "for", "func",
    "go", "goto", "if", "import", "interface", "map", "package", "range", "return", "select", "struct",
    "switch", "type", "var", "nil", "true", "false",
  ],
  js: [
    "async", "await", "break", "case", "catch", "class", "const", "continue", "default", "else", "export",
    "extends", "finally", "for", "from", "function", "if", "import", "interface", "let", "new", "null",
    "return", "switch", "throw", "true", "false", "try", "type", "undefined", "var", "while",
  ],
  python: [
    "and", "as", "assert", "async", "await", "class", "def", "elif", "else", "except", "finally", "for",
    "from", "if", "import", "in", "is", "lambda", "None", "not", "or", "pass", "raise", "return", "try",
    "True", "False", "while", "with", "yield",
  ],
  sql: [
    "select", "from", "where", "insert", "into", "values", "update", "set", "delete", "create", "table",
    "alter", "drop", "join", "left", "right", "inner", "on", "group", "by", "order", "limit", "offset",
    "and", "or", "not", "null", "as", "distinct", "primary", "key", "references", "default", "index",
  ],
  bash: ["if", "then", "else", "fi", "for", "do", "done", "while", "case", "esac", "function", "echo", "export", "cd"],
};

const ALIAS: Record<string, keyof typeof KEYWORDS> = {
  golang: "go",
  javascript: "js",
  typescript: "js",
  ts: "js",
  tsx: "js",
  jsx: "js",
  py: "python",
  postgresql: "sql",
  psql: "sql",
  sh: "bash",
  shell: "bash",
  powershell: "bash",
};

type Familia = "go" | "js" | "python" | "sql" | "bash" | null;

function familia(lang?: string): Familia {
  if (!lang) return null;
  const l = lang.toLowerCase();
  const chave = (ALIAS[l] ?? l) as keyof typeof KEYWORDS;
  return chave in KEYWORDS ? (chave as Familia) : null;
}

function comentarioRegex(fam: Familia): string {
  switch (fam) {
    case "python":
    case "bash":
      return "#[^\\n]*";
    case "sql":
      return "--[^\\n]*|/\\*[\\s\\S]*?\\*/";
    default:
      return "//[^\\n]*|/\\*[\\s\\S]*?\\*/";
  }
}

export function tokenizar(code: string, lang?: string): Token[] {
  const fam = familia(lang);
  if (!fam) return [{ type: "plain", text: code }];

  const backtick = fam === "go" || fam === "js" ? "|`[^`]*`" : "";
  const re = new RegExp(
    [
      `(?<comment>${comentarioRegex(fam)})`,
      `(?<string>"(?:\\\\.|[^"\\\\\\n])*"|'(?:\\\\.|[^'\\\\\\n])*'${backtick})`,
      "(?<number>\\b\\d+(?:\\.\\d+)?\\b)",
      "(?<word>[A-Za-z_][A-Za-z0-9_]*)",
    ].join("|"),
    "g",
  );
  const palavras = new Set(KEYWORDS[fam].map((k) => (fam === "sql" ? k.toLowerCase() : k)));

  const tokens: Token[] = [];
  let ultimo = 0;
  const push = (type: TokenType, text: string) => {
    if (!text) return;
    const prev = tokens[tokens.length - 1];
    if (prev && prev.type === type) prev.text += text;
    else tokens.push({ type, text });
  };

  for (const m of Array.from(code.matchAll(re))) {
    const idx = m.index ?? 0;
    push("plain", code.slice(ultimo, idx));
    const g = m.groups ?? {};
    if (g.comment !== undefined) push("comment", g.comment);
    else if (g.string !== undefined) push("string", g.string);
    else if (g.number !== undefined) push("number", g.number);
    else if (g.word !== undefined) {
      const w = fam === "sql" ? g.word.toLowerCase() : g.word;
      push(palavras.has(w) ? "keyword" : "plain", g.word);
    }
    ultimo = idx + m[0].length;
  }
  push("plain", code.slice(ultimo));
  return tokens;
}
