import { describe, expect, it } from "vitest";

import { destinoSeguro } from "@/lib/navigation";

describe("destinoSeguro", () => {
  it("aceita caminhos internos", () => {
    expect(destinoSeguro("/artigos/go-intro")).toBe("/artigos/go-intro");
    expect(destinoSeguro("/busca?q=go")).toBe("/busca?q=go");
  });

  it("usa o padrão quando vazio ou ausente", () => {
    expect(destinoSeguro(null)).toBe("/");
    expect(destinoSeguro(undefined)).toBe("/");
    expect(destinoSeguro("")).toBe("/");
    expect(destinoSeguro("", "/inicio")).toBe("/inicio");
  });

  it("bloqueia open redirect", () => {
    expect(destinoSeguro("//evil.com")).toBe("/");
    expect(destinoSeguro("/\\evil.com")).toBe("/");
    expect(destinoSeguro("https://evil.com")).toBe("/");
    expect(destinoSeguro("javascript:alert(1)")).toBe("/");
    expect(destinoSeguro("/\n/evil.com")).toBe("/");
  });

  it("não devolve para login/registro", () => {
    expect(destinoSeguro("/login")).toBe("/");
    expect(destinoSeguro("/registro?next=/x")).toBe("/");
    expect(destinoSeguro("/loginx")).toBe("/loginx");
  });
});
