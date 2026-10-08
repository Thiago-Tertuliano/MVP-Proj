import { z } from "zod";

/** Regras espelham o backend (RegistrarRequest/LoginRequest) para falhar cedo, sem round-trip. */

const email = z
  .string()
  .trim()
  .min(1, "Informe seu e-mail.")
  .max(255, "O e-mail pode ter no máximo 255 caracteres.")
  .email("Digite um e-mail válido, como nome@exemplo.com.");

export const loginSchema = z.object({
  email,
  senha: z.string().min(1, "Informe sua senha."),
});

export const registroSchema = z.object({
  nome: z.string().trim().min(2, "Informe seu nome (mínimo de 2 letras).").max(150, "O nome pode ter no máximo 150 caracteres."),
  email,
  senha: z
    .string()
    .min(8, "A senha precisa ter pelo menos 8 caracteres.")
    .max(72, "A senha pode ter no máximo 72 caracteres."),
});

export type LoginValues = z.infer<typeof loginSchema>;
export type RegistroValues = z.infer<typeof registroSchema>;
