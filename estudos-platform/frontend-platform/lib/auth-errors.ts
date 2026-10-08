import { ApiError } from "@/lib/api";

export type FormAlertData = {
  tone: "danger" | "warning";
  titulo: string;
  mensagem: string;
};

const GENERICO: FormAlertData = {
  tone: "danger",
  titulo: "Algo deu errado",
  mensagem: "Não conseguimos concluir agora. Tente novamente em instantes.",
};

function comuns(err: unknown): FormAlertData | null {
  if (!(err instanceof ApiError)) return GENERICO;
  if (err.isRateLimited) {
    return {
      tone: "warning",
      titulo: "Muitas tentativas",
      mensagem: "Por segurança, aguarde 1 minuto antes de tentar de novo.",
    };
  }
  if (err.isNetwork) {
    return {
      tone: "danger",
      titulo: "Sem conexão com o servidor",
      mensagem: "Verifique sua internet e tente novamente.",
    };
  }
  if (err.status >= 500) {
    return {
      tone: "danger",
      titulo: "Servidor indisponível",
      mensagem: "Tivemos um problema do nosso lado. Tente novamente em instantes.",
    };
  }
  return null;
}

/** Copy amigável para POST /auth/login. Nunca expõe stack nem mensagem crua da API. */
export function erroLogin(err: unknown): FormAlertData {
  const base = comuns(err);
  if (base) return base;
  const status = (err as ApiError).status;
  if (status === 401) {
    return { tone: "danger", titulo: "Não foi possível entrar", mensagem: "E-mail ou senha incorretos." };
  }
  if (status === 400 || status === 422) {
    return { tone: "danger", titulo: "Dados inválidos", mensagem: "Confira o e-mail e a senha informados." };
  }
  return GENERICO;
}

/** Copy amigável para POST /auth/registrar. */
export function erroRegistro(err: unknown): FormAlertData {
  const base = comuns(err);
  if (base) return base;
  const status = (err as ApiError).status;
  if (status === 409) {
    return {
      tone: "danger",
      titulo: "E-mail já cadastrado",
      mensagem: "Já existe uma conta com esse e-mail. Tente entrar ou use outro e-mail.",
    };
  }
  if (status === 400 || status === 422) {
    return {
      tone: "danger",
      titulo: "Dados inválidos",
      mensagem: "Confira nome, e-mail e senha (mínimo de 8 caracteres).",
    };
  }
  return GENERICO;
}
