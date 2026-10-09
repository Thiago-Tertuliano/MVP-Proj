/**
 * Cliente HTTP do Relp.
 *
 * - Base URL vem de `NEXT_PUBLIC_API_URL` (no servidor, `API_URL` tem prioridade — útil em Docker).
 * - `credentials: "include"`: o JWT vive em cookie HttpOnly; o front nunca guarda token.
 * - Erros viram `ApiError` com copy amigável: nada de stack/JSON cru chega na UI.
 * - 401 fora das rotas públicas de auth: renova a sessão em silêncio (1 refresh por vez) e refaz a chamada.
 */

const DEFAULT_API_URL = "http://localhost:8080";

export function apiBaseUrl(): string {
  const isServer = typeof window === "undefined";
  const raw =
    (isServer ? process.env.API_URL : undefined) || process.env.NEXT_PUBLIC_API_URL || DEFAULT_API_URL;
  return `${raw.replace(/\/+$/, "")}/api/v1`;
}

/* ------------------------------------------------------------------ erros */

export class ApiError extends Error {
  /** Status HTTP; `0` quando não houve resposta (rede fora, CORS, API desligada). */
  readonly status: number;
  /** Mensagem de domínio devolvida pela API (`{"erro": "..."}`), só para uso interno. */
  readonly serverMessage?: string;

  constructor(status: number, serverMessage?: string) {
    super(mensagemPorStatus(status, serverMessage));
    this.name = "ApiError";
    this.status = status;
    this.serverMessage = serverMessage;
  }

  get isNetwork() {
    return this.status === 0;
  }
  get isUnauthorized() {
    return this.status === 401;
  }
  get isNotFound() {
    return this.status === 404;
  }
  get isRateLimited() {
    return this.status === 429;
  }
}

export function mensagemPorStatus(status: number, serverMessage?: string): string {
  if (status === 0) return "Não foi possível conectar ao servidor. Verifique sua conexão e tente novamente.";
  if (status === 401) return "Sua sessão expirou ou é inválida. Entre novamente.";
  if (status === 403) return "Você não tem permissão para fazer isso.";
  if (status === 404) return "Não encontramos o que você procurava.";
  // 409 carrega mensagem de domínio já em português (ex.: "e-mail já cadastrado").
  if (status === 409) return serverMessage || "Essa ação conflita com o estado atual. Atualize a página e tente de novo.";
  if (status === 429) return "Muitas tentativas em pouco tempo. Aguarde 1 minuto e tente novamente.";
  if (status === 400 || status === 422) return "Os dados enviados não são válidos. Confira e tente novamente.";
  if (status >= 500) return "O servidor teve um problema. Tente novamente em instantes.";
  return "Algo deu errado. Tente novamente.";
}

/* ------------------------------------------------------------------ sessão (refresh) */

/** Rotas públicas de auth: 401 aqui é "credencial errada", não "sessão vencida". */
const ROTAS_SEM_REFRESH = ["/auth/login", "/auth/registrar", "/auth/refresh"];

type Listener = () => void;
const sessionExpiredListeners = new Set<Listener>();

/** Avisa quando o refresh falha — a camada de sessão derruba o usuário para visitante. */
export function onSessionExpired(listener: Listener): () => void {
  sessionExpiredListeners.add(listener);
  return () => sessionExpiredListeners.delete(listener);
}

let refreshEmAndamento: Promise<boolean> | null = null;

/** Single-flight: várias chamadas 401 simultâneas compartilham um único POST /auth/refresh. */
export function renovarSessao(): Promise<boolean> {
  if (!refreshEmAndamento) {
    refreshEmAndamento = send("/auth/refresh", { method: "POST" })
      .then((res) => res.ok)
      .catch(() => false)
      .finally(() => {
        refreshEmAndamento = null;
      });
  }
  return refreshEmAndamento;
}

/* ------------------------------------------------------------------ requisição */

export type ApiOptions = {
  method?: "GET" | "POST" | "PUT" | "DELETE";
  body?: unknown;
  query?: Record<string, string | number | undefined>;
  signal?: AbortSignal;
  /** Passa direto ao fetch (Next usa `no-store` nas páginas de servidor). */
  cache?: RequestCache;
  /** Liga/desliga o refresh silencioso. Padrão: só no browser. */
  refresh?: boolean;
};

function buildUrl(path: string, query?: ApiOptions["query"]): string {
  // Em produção sem domínio próprio a base é relativa ("/") e o Next faz proxy para a API.
  const origem = typeof window === "undefined" ? undefined : window.location.origin;
  const url = new URL(`${apiBaseUrl()}${path}`, origem);
  if (query) {
    for (const [k, v] of Object.entries(query)) {
      if (v !== undefined && v !== "") url.searchParams.set(k, String(v));
    }
  }
  return url.toString();
}

async function send(path: string, opts: ApiOptions): Promise<Response> {
  const headers: Record<string, string> = { Accept: "application/json" };
  let body: string | undefined;
  if (opts.body !== undefined) {
    headers["Content-Type"] = "application/json";
    body = JSON.stringify(opts.body);
  }
  try {
    return await fetch(buildUrl(path, opts.query), {
      method: opts.method ?? "GET",
      headers,
      body,
      credentials: "include",
      signal: opts.signal,
      cache: opts.cache,
    });
  } catch (err) {
    if (err instanceof DOMException && err.name === "AbortError") throw err;
    throw new ApiError(0);
  }
}

async function parse<T>(res: Response): Promise<T> {
  if (res.status === 204) return undefined as T;
  const isJson = res.headers.get("content-type")?.includes("application/json");
  const data: unknown = isJson ? await res.json().catch(() => undefined) : undefined;
  if (!res.ok) {
    const msg =
      data && typeof data === "object" && "erro" in data && typeof (data as { erro: unknown }).erro === "string"
        ? (data as { erro: string }).erro
        : undefined;
    throw new ApiError(res.status, msg);
  }
  return data as T;
}

export async function api<T>(path: string, opts: ApiOptions = {}): Promise<T> {
  const res = await send(path, opts);

  const podeRenovar =
    res.status === 401 && (opts.refresh ?? typeof window !== "undefined") && !ROTAS_SEM_REFRESH.includes(path);

  if (podeRenovar) {
    if (await renovarSessao()) {
      return parse<T>(await send(path, opts));
    }
    sessionExpiredListeners.forEach((fn) => fn());
  }
  return parse<T>(res);
}
