import { describe, expect, it } from "vitest";

import { tokenizar } from "@/lib/highlight";

const juntar = (tokens: { text: string }[]) => tokens.map((t) => t.text).join("");

describe("tokenizar", () => {
  it("nunca perde nem altera texto (round-trip)", () => {
    const codigo = 'package main\n\n// olá\nfunc main() {\n\tfmt.Println("oi", 42)\n}\n';
    expect(juntar(tokenizar(codigo, "go"))).toBe(codigo);
  });

  it("go: palavras-chave, strings, números e comentários", () => {
    const tokens = tokenizar('func x() { return "a" } // fim 10', "go");
    const tipos = Object.fromEntries(tokens.filter((t) => t.type !== "plain").map((t) => [t.text.trim(), t.type]));
    expect(tipos["func"]).toBe("keyword");
    expect(tipos["return"]).toBe("keyword");
    expect(tipos['"a"']).toBe("string");
    expect(tipos["// fim 10"]).toBe("comment"); // número dentro de comentário não vira number
  });

  it("não confunde // dentro de string com comentário", () => {
    const tokens = tokenizar('url := "http://x.dev"', "go");
    expect(tokens.find((t) => t.type === "comment")).toBeUndefined();
    expect(tokens.find((t) => t.type === "string")?.text).toBe('"http://x.dev"');
  });

  it("aceita apelidos de linguagem (ts, py, sh)", () => {
    expect(tokenizar("const a = 1", "ts").some((t) => t.type === "keyword")).toBe(true);
    expect(tokenizar("def f(): pass  # x", "py").some((t) => t.type === "comment")).toBe(true);
    expect(tokenizar("echo ok # x", "sh").some((t) => t.type === "comment")).toBe(true);
  });

  it("sql ignora maiúsculas nas palavras-chave", () => {
    const kw = tokenizar("SELECT * FROM t", "sql").filter((t) => t.type === "keyword");
    expect(kw.map((t) => t.text)).toEqual(["SELECT", "FROM"]);
  });

  it("linguagem desconhecida ou ausente devolve texto simples", () => {
    expect(tokenizar("qualquer coisa", "brainfuck")).toEqual([{ type: "plain", text: "qualquer coisa" }]);
    expect(tokenizar("x", undefined)).toEqual([{ type: "plain", text: "x" }]);
  });
});
