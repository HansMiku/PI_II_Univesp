# Supabase Migration — LegisAnalytica Gold

Esta pasta contém o(s) script(s) SQL de criação da estrutura relacional Gold para o Supabase/PostgreSQL, idêntica ao modelo de dados do Databricks Gold Layer.

## Arquivos

- `001_initial_schema_supabase.sql` — Migration SQL: cria as 26 tabelas (ordem topológica em 6 níveis), 60 FKs, 30 índices, 2 constraints UNIQUE (`lote`), 26 comentários e os tipos nativos do PostgreSQL (TEXT, NUMERIC, etc.), além das regras de dependência das FKs.
- `002_add_autoria_principal.sql` — Adiciona a coluna `autoria.principal`, que diferencia o autor principal dos demais signatários/apoiadores.

## Uso

1. Abra o painel SQL do Supabase e **execute este script (`001_initial_schema_supabase.sql`)** no schema `public` vazio (ou desejado) antes de qualquer exportação de dados.
2. Se o schema já existir, execute também `002_add_autoria_principal.sql`.
3. Nunca execute este SQL no Databricks/Spark SQL — é exclusivo para PostgreSQL.
4. Caso precise reinstalar, primeiro delete as tabelas manualmente (ou descomente os blocos de DROP indicados no topo do SQL).

### Base já existente

Depois de aplicar a `002`, preencha a coluna `principal` na base atual:

```powershell
cd 05_Aplicacao\backend
go run ./cmd/backfill-autoria
```

O comando:
- marca como principal o primeiro signatário da Câmara (`ordem_assinatura = 1`);
- consulta a API de autorias do Senado e marca como principal apenas quem tem `IndicadorAutorPrincipal = "Sim"`;
- é idempotente e pode ser executado novamente.

Também é possível refazer a carga Gold/exportação: o notebook `03_Gold/LegisAnalytica - Gold Layer DDL.ipynb` já grava a coluna `principal` no schema Gold.

**Importante:**
- A exportação dos dados só deve ser feita após a criação correta dessa estrutura!
- O script aqui é referência única para criação/reinstalação do modelo LegisAnalytica Gold fora do Databricks.
- A aplicação é manual, no SQL Editor do Supabase; este repositório versiona o script como fonte de referência (o notebook de exportação o aplica automaticamente apenas como fallback, em banco vazio).
