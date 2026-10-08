import { describe, expect, it } from "vitest";

import { ApiError } from "@/lib/api";
import { erroLogin, erroRegistro } from "@/lib/auth-errors";
import { loginSchema, registroSchema } from "@/lib/validation";

function mensagens(result: { success: boolean; error?: { issues: { path: PropertyKey[]; message: string }[] } }) {
  // Como o react-hook-form: vale a primeira mensagem de cada campo.
  const out: Record<string, string> = {};
  for (const i of result.error?.issues ?? []) out[String(i.path[0])] ??= i.message;
  return out;
}

describe("validação de formulários (Zod)", () => {
  it("login: campos vazios geram mensagens", () => {
    const r = loginSchema.safeParse({ email: "", senha: "" });
    expect(r.success).toBe(false);
    expect(mensagens(r)).toEqual({ email: "Informe seu e-mail.", senha: "Informe sua senha." });
  });

  it("login: formato de e-mail inválido", () => {
    const r = loginSchema.safeParse({ email: "sem-arroba", senha: "x" });
    expect(mensagens(r).email).toMatch(/e-mail válido/);
  });

  it("login: aceita e-mail com espaços nas pontas (trim)", () => {
    const r = loginSchema.safeParse({ email: "  ana@exemplo.com ", senha: "x" });
    expect(r.success).toBe(true);
    expect(r.data?.email).toBe("ana@exemplo.com");
  });

  it("registro: senha precisa de 8+ caracteres (regra do backend)", () => {
    const r = registroSchema.safeParse({ nome: "Ana", email: "ana@exemplo.com", senha: "1234567" });
    expect(mensagens(r).senha).toMatch(/8 caracteres/);
    expect(registroSchema.safeParse({ nome: "Ana", email: "ana@exemplo.com", senha: "12345678" }).success).toBe(true);
  });

  it("registro: nome obrigatório", () => {
    const r = registroSchema.safeParse({ nome: " ", email: "ana@exemplo.com", senha: "12345678" });
    expect(mensagens(r).nome).toMatch(/nome/i);
  });
});

describe("copy de erro de auth", () => {
  it("login 401 = credenciais incorretas, sem detalhes técnicos", () => {
    const r = erroLogin(new ApiError(401, "credenciais inválidas"));
    expect(r.tone).toBe("danger");
    expect(r.mensagem).toBe("E-mail ou senha incorretos.");
  });

  it("login 429 = aviso de rate limit", () => {
    const r = erroLogin(new ApiError(429));
    expect(r.tone).toBe("warning");
    expect(r.mensagem).toMatch(/1 minuto/);
  });

  it("registro 409 = e-mail já cadastrado", () => {
    expect(erroRegistro(new ApiError(409, "e-mail já cadastrado")).titulo).toMatch(/já cadastrado/i);
  });

  it("erros inesperados caem em copy genérica (nunca vazam Error.message)", () => {
    const r = erroLogin(new Error("TypeError: x is undefined at foo.js:12"));
    expect(JSON.stringify(r)).not.toMatch(/TypeError|foo\.js/);
  });

  it("5xx e rede têm copy própria", () => {
    expect(erroLogin(new ApiError(503)).titulo).toMatch(/indispon/i);
    expect(erroLogin(new ApiError(0)).titulo).toMatch(/conex/i);
  });
});
