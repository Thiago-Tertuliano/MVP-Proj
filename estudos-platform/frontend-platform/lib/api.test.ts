import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { api, ApiError, apiBaseUrl, onSessionExpired } from "@/lib/api";

function json(status: number, body?: unknown): Response {
  return new Response(body === undefined ? null : JSON.stringify(body), {
    status,
    headers: body === undefined ? {} : { "content-type": "application/json" },
  });
}

/** Espera a promise rejeitar e devolve o erro já tipado. */
async function falha(p: Promise<unknown>): Promise<ApiError> {
  try {
    await p;
  } catch (e) {
    return e as ApiError;
  }
  throw new Error("a chamada deveria ter falhado");
}

describe("api client", () => {
  const fetchMock = vi.fn<typeof fetch>();

  beforeEach(() => {
    fetchMock.mockReset();
    vi.stubGlobal("fetch", fetchMock);
  });
  afterEach(() => {
    vi.unstubAllGlobals();
    vi.unstubAllEnvs();
  });

  it("usa NEXT_PUBLIC_API_URL + /api/v1 e envia cookies (credentials: include)", async () => {
    vi.stubEnv("API_URL", "");
    vi.stubEnv("NEXT_PUBLIC_API_URL", "https://api.relp.test/");
    fetchMock.mockResolvedValue(json(200, { ok: true }));

    await api("/trilhas", { query: { limit: 10, vazio: undefined } });

    const [url, init] = fetchMock.mock.calls[0];
    expect(String(url)).toBe("https://api.relp.test/api/v1/trilhas?limit=10");
    expect(init?.credentials).toBe("include");
    expect(apiBaseUrl()).toBe("https://api.relp.test/api/v1");
  });

  it("serializa body em JSON e define Content-Type só quando há body", async () => {
    fetchMock.mockResolvedValue(json(200, {}));
    await api("/auth/login", { method: "POST", body: { email: "a@b.co", senha: "x" } });
    await api("/trilhas");

    const comBody = fetchMock.mock.calls[0][1];
    const semBody = fetchMock.mock.calls[1][1];
    expect((comBody?.headers as Record<string, string>)["Content-Type"]).toBe("application/json");
    expect(comBody?.body).toBe(JSON.stringify({ email: "a@b.co", senha: "x" }));
    expect((semBody?.headers as Record<string, string>)["Content-Type"]).toBeUndefined();
  });

  it("padroniza erros HTTP em ApiError com copy amigável (sem JSON cru)", async () => {
    fetchMock.mockResolvedValue(json(422, { erro: "Key: 'RegistrarRequest.Email' Error:Field validation" }));
    const err = await falha(api("/auth/registrar", { method: "POST" }));

    expect(err).toBeInstanceOf(ApiError);
    expect(err.status).toBe(422);
    expect(err.message).not.toContain("Key:");
    expect(err.message).toMatch(/dados enviados/i);
  });

  it("429 vira mensagem de rate limit; 409 repassa a mensagem de domínio", async () => {
    fetchMock.mockResolvedValueOnce(json(429, { erro: "muitas tentativas" }));
    const rate = await falha(api("/auth/login", { method: "POST" }));
    expect(rate.isRateLimited).toBe(true);
    expect(rate.message).toMatch(/1 minuto/);

    fetchMock.mockResolvedValueOnce(json(409, { erro: "e-mail já cadastrado" }));
    const conflito = await falha(api("/auth/registrar", { method: "POST" }));
    expect(conflito.message).toBe("e-mail já cadastrado");
  });

  it("falha de rede vira ApiError status 0", async () => {
    fetchMock.mockRejectedValue(new TypeError("Failed to fetch"));
    const err = await falha(api("/trilhas"));
    expect(err).toBeInstanceOf(ApiError);
    expect(err.isNetwork).toBe(true);
    expect(err.message).not.toContain("Failed to fetch");
  });

  describe("refresh silencioso em 401", () => {
    it("renova a sessão e refaz a chamada original", async () => {
      fetchMock
        .mockResolvedValueOnce(json(401, { erro: "token inválido" })) // chamada original
        .mockResolvedValueOnce(json(200, {})) // POST /auth/refresh
        .mockResolvedValueOnce(json(200, { artigo_id: "a1", concluido: true })); // retry

      const resp = await api<{ concluido: boolean }>("/progresso/artigos/a1", {
        method: "PUT",
        body: { concluido: true },
        refresh: true,
      });

      expect(resp.concluido).toBe(true);
      expect(fetchMock).toHaveBeenCalledTimes(3);
      expect(String(fetchMock.mock.calls[1][0])).toContain("/auth/refresh");
      expect(fetchMock.mock.calls[1][1]?.method).toBe("POST");
      expect(String(fetchMock.mock.calls[2][0])).toContain("/progresso/artigos/a1");
    });

    it("várias chamadas 401 simultâneas compartilham UM refresh", async () => {
      let refreshes = 0;
      fetchMock.mockImplementation(async (input) => {
        const url = String(input);
        if (url.includes("/auth/refresh")) {
          refreshes++;
          await new Promise((r) => setTimeout(r, 10));
          return json(200, {});
        }
        // 1ª rodada 401; depois do refresh responde 200
        return refreshes === 0 ? json(401, { erro: "expirado" }) : json(200, { ok: true });
      });

      await Promise.all([
        api("/progresso/trilhas/1", { refresh: true }),
        api("/progresso/trilhas/2", { refresh: true }),
        api("/artigos/x/anotacoes", { refresh: true }),
      ]);
      expect(refreshes).toBe(1);
    });

    it("se o refresh falha, avisa a sessão e propaga o 401", async () => {
      const expirou = vi.fn();
      const off = onSessionExpired(expirou);
      fetchMock
        .mockResolvedValueOnce(json(401, { erro: "expirado" }))
        .mockResolvedValueOnce(json(401, { erro: "refresh inválido" }));

      const err = await falha(api("/progresso/trilhas/1", { refresh: true }));
      off();

      expect(err.isUnauthorized).toBe(true);
      expect(expirou).toHaveBeenCalledTimes(1);
      expect(fetchMock).toHaveBeenCalledTimes(2); // sem laço infinito
    });

    it.each(["/auth/login", "/auth/registrar", "/auth/refresh"])(
      "não tenta refresh em %s (401 = credencial errada)",
      async (path) => {
        fetchMock.mockResolvedValue(json(401, { erro: "credenciais inválidas" }));
        const err = await falha(api(path, { method: "POST", refresh: true }));
        expect(err.isUnauthorized).toBe(true);
        expect(fetchMock).toHaveBeenCalledTimes(1);
      },
    );

    it("GET /auth/me e POST /auth/logout PODEM renovar (precisam do access token)", async () => {
      fetchMock
        .mockResolvedValueOnce(json(401, {}))
        .mockResolvedValueOnce(json(200, {}))
        .mockResolvedValueOnce(json(200, { id: "u1", nome: "Ana", email: "ana@x.co" }));
      const me = await api<{ nome: string }>("/auth/me", { refresh: true });
      expect(me.nome).toBe("Ana");
    });

    it("sem refresh habilitado (servidor), 401 só propaga", async () => {
      fetchMock.mockResolvedValue(json(401, {}));
      const err = await falha(api("/progresso/trilhas/1", { refresh: false }));
      expect(err.isUnauthorized).toBe(true);
      expect(fetchMock).toHaveBeenCalledTimes(1);
    });
  });
});
