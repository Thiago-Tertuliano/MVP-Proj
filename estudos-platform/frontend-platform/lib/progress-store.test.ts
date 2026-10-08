import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import {
  aplicarMarcacao,
  carregarProgresso,
  getProgressoEntry,
  limparProgresso,
  marcarArtigo,
} from "@/lib/progress-store";
import type { ProgressoTrilha } from "@/lib/types";

const base: ProgressoTrilha = {
  trilha_id: "t1",
  concluidos: 1,
  total: 3,
  percentual: 100 / 3,
  artigos_concluidos: ["a1"],
};

function resposta(status: number, body: unknown = {}) {
  return new Response(JSON.stringify(body), { status, headers: { "content-type": "application/json" } });
}

describe("aplicarMarcacao", () => {
  it("marcar sobe contador e percentual", () => {
    const r = aplicarMarcacao(base, "a2", true);
    expect(r.concluidos).toBe(2);
    expect(r.artigos_concluidos).toEqual(["a1", "a2"]);
    expect(r.percentual).toBeCloseTo(66.67, 1);
  });
  it("é idempotente e desmarcar volta atrás", () => {
    expect(aplicarMarcacao(base, "a1", true).concluidos).toBe(1);
    const r = aplicarMarcacao(base, "a1", false);
    expect(r.concluidos).toBe(0);
    expect(r.percentual).toBe(0);
  });
  it("trilha vazia não divide por zero", () => {
    expect(aplicarMarcacao({ ...base, total: 0, concluidos: 0, artigos_concluidos: [] }, "a1", true).percentual).toBe(0);
  });
});

describe("store de progresso (optimistic UI)", () => {
  const fetchMock = vi.fn<typeof fetch>();

  beforeEach(() => {
    limparProgresso();
    fetchMock.mockReset();
    vi.stubGlobal("fetch", fetchMock);
  });
  afterEach(() => vi.unstubAllGlobals());

  it("carrega e cacheia: segunda chamada não vai à rede", async () => {
    fetchMock.mockResolvedValue(resposta(200, base));
    await carregarProgresso("t1");
    await carregarProgresso("t1");
    expect(fetchMock).toHaveBeenCalledTimes(1);
    expect(getProgressoEntry("t1")?.status).toBe("ready");
    expect(getProgressoEntry("t1")?.data?.concluidos).toBe(1);
  });

  it("marca como lido ANTES da API responder", async () => {
    fetchMock.mockResolvedValueOnce(resposta(200, base));
    await carregarProgresso("t1");

    let liberar!: (r: Response) => void;
    fetchMock.mockReturnValueOnce(new Promise<Response>((res) => (liberar = res)));

    const pendente = marcarArtigo({ artigoId: "a2", trilhaId: "t1", concluido: true });
    // API ainda não respondeu e o cache já reflete o novo estado
    expect(getProgressoEntry("t1")?.data?.artigos_concluidos).toContain("a2");
    expect(getProgressoEntry("t1")?.data?.concluidos).toBe(2);

    liberar(resposta(200, { artigo_id: "a2", concluido: true }));
    await pendente;
    expect(getProgressoEntry("t1")?.data?.concluidos).toBe(2);
  });

  it("se o PUT falha, desfaz o optimistic e propaga o erro", async () => {
    fetchMock.mockResolvedValueOnce(resposta(200, base));
    await carregarProgresso("t1");
    fetchMock.mockResolvedValueOnce(resposta(500, { erro: "erro interno" }));

    await expect(marcarArtigo({ artigoId: "a2", trilhaId: "t1", concluido: true })).rejects.toMatchObject({
      status: 500,
    });
    expect(getProgressoEntry("t1")?.data?.artigos_concluidos).toEqual(["a1"]);
    expect(getProgressoEntry("t1")?.data?.concluidos).toBe(1);
  });

  it("limparProgresso (logout) esvazia o cache", async () => {
    fetchMock.mockResolvedValue(resposta(200, base));
    await carregarProgresso("t1");
    limparProgresso();
    expect(getProgressoEntry("t1")).toBeUndefined();
  });

  it("erro de carregamento fica em 'error' e permite nova tentativa", async () => {
    fetchMock.mockResolvedValueOnce(resposta(500, { erro: "x" }));
    await carregarProgresso("t1");
    expect(getProgressoEntry("t1")?.status).toBe("error");

    fetchMock.mockResolvedValueOnce(resposta(200, base));
    await carregarProgresso("t1");
    expect(getProgressoEntry("t1")?.status).toBe("ready");
  });
});
