-- ============================================================
-- LegisData — Migration 002
-- Adiciona a distinção entre autor principal e demais signatários
-- na tabela public.autoria.
--
-- Contexto:
--   No Senado, a API de autorias retorna todos os senadores que
--   assinaram/apoiaram uma matéria. O campo IndicadorAutorPrincipal
--   identifica quem é o autor principal; os demais são apoiadores.
--
--   Na Câmara, o primeiro signatário (ordem_assinatura = 1) é o
--   autor principal.
--
-- Execute no SQL Editor do Supabase.
-- ============================================================

ALTER TABLE public.autoria
    ADD COLUMN IF NOT EXISTS principal BOOLEAN NOT NULL DEFAULT false;

-- Câmara: primeiro signatário é o autor principal.
UPDATE public.autoria
SET principal = (ordem_assinatura = 1)
WHERE casa = 'CAMARA';

-- Senado:
-- O valor será preenchido pela carga Gold (campo IndicadorAutorPrincipal).
-- Para a base já existente, execute o backfill descrito em
-- 03_Gold/supabase/migrations/README.md antes de usar o filtro de autoria.
CREATE INDEX IF NOT EXISTS idx_autoria_autor_principal
    ON public.autoria (casa, id_autor, principal);
