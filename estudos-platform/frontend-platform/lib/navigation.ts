/**
 * Destino pós-login vindo de `?next=`. Só aceita caminho interno ("/artigos/x"); qualquer coisa que
 * possa sair do site (//evil.com, /\evil.com, https://…, javascript:) cai na home — evita open redirect.
 */
export function destinoSeguro(next: string | null | undefined, padrao = "/"): string {
  if (!next) return padrao;
  if (!next.startsWith("/")) return padrao;
  if (next.startsWith("//") || next.includes("\\")) return padrao;
  // Caracteres de controle (quebra de linha, tab) podem enganar o parser de URL do navegador.
  // eslint-disable-next-line no-control-regex
  if (/[\u0000-\u001f]/.test(next)) return padrao;
  // Não volta para as telas de autenticação (loop).
  if (/^\/(login|registro)(\/|\?|$)/.test(next)) return padrao;
  return next;
}
