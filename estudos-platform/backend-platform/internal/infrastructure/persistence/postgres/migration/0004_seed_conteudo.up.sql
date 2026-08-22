-- autor de demo. senha = senha1234
INSERT INTO usuarios (id, nome, email, senha_hash, status)
VALUES (
  '11111111-1111-1111-1111-111111111111',
  'Autor Seed',
  'autor.seed@estudos.local',
  '$2a$10$GnWGzZlXrOg9M.dqhPxVR.qVwdW1.MKsO7lVtsv.iupzECQ089Em2',
  'ativa'
) ON CONFLICT (email) DO NOTHING;

INSERT INTO trilhas (id, slug, titulo, descricao, ordem, publicada, created_at, updated_at)
VALUES (
  '22222222-2222-2222-2222-222222222222',
  'go-basico',
  'Go Básico',
  'Fundamentos de Go para a plataforma',
  0, true, now(), now()
) ON CONFLICT (slug) DO NOTHING;

INSERT INTO modulos (id, trilha_id, slug, titulo, descricao, ordem) VALUES
  ('33333333-3333-3333-3333-333333333333', '22222222-2222-2222-2222-222222222222', 'sintaxe', 'Sintaxe', 'Pacotes e tipos', 0),
  ('44444444-4444-4444-4444-444444444444', '22222222-2222-2222-2222-222222222222', 'interfaces', 'Interfaces', 'Contratos', 1)
ON CONFLICT DO NOTHING;

INSERT INTO artigos (id, slug, titulo, conteudo, metadados, autor_id, status, publicado_em)
VALUES
  ('55555555-5555-5555-5555-555555555555', 'pacotes-em-go', 'Pacotes em Go',
   '{"blocks":[{"type":"p","text":"Todo arquivo declara um package."}]}', '{}',
   '11111111-1111-1111-1111-111111111111', 'publicado', now()),
  ('66666666-6666-6666-6666-666666666666', 'structs-e-metodos', 'Structs e métodos',
   '{"blocks":[{"type":"p","text":"Receiver por valor vs ponteiro."}]}', '{}',
   '11111111-1111-1111-1111-111111111111', 'publicado', now()),
  ('77777777-7777-7777-7777-777777777777', 'interfaces-implicitas', 'Interfaces implícitas',
   '{"blocks":[{"type":"p","text":"Satisfaz sem implements."}]}', '{}',
   '11111111-1111-1111-1111-111111111111', 'publicado', now())
ON CONFLICT (slug) DO NOTHING;
