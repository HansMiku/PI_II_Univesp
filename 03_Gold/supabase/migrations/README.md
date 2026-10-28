# Supabase Migration — LegisAnalytica Gold

Esta pasta contém o(s) script(s) SQL de criação da estrutura relacional Gold para o Supabase/PostgreSQL, idêntica ao modelo de dados do Databricks Gold Layer.

## Arquivos

- `001_initial_schema_supabase.sql` — Migration SQL: cria as 26 tabelas (ordem topológica em 6 níveis), 60 FKs, 30 índices, 2 constraints UNIQUE (`lote`), 26 comentários e os tipos nativos do PostgreSQL (TEXT, NUMERIC, etc.), além das regras de dependência das FKs.

## Uso

1. Abra o painel SQL do Supabase e **execute este script (`001_initial_schema_supabase.sql`)** no schema `public` vazio (ou desejado) antes de qualquer exportação de dados.
2. Nunca execute este SQL no Databricks/Spark SQL — é exclusivo para PostgreSQL.
3. Caso precise reinstalar, primeiro delete as tabelas manualmente (ou descomente os blocos de DROP indicados no topo do SQL).

**Importante:**
- A exportação dos dados só deve ser feita após a criação correta dessa estrutura!
- O script aqui é referência única para criação/reinstalação do modelo LegisAnalytica Gold fora do Databricks.
- A aplicação é manual, no SQL Editor do Supabase; este repositório versiona o script como fonte de referência (o notebook de exportação o aplica automaticamente apenas como fallback, em banco vazio).
