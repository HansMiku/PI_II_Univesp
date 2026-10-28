# Gold — Data Warehouse

Responsável: Bruce

Esta pasta contém os códigos relacionados à modelagem e organização dos dados consolidados utilizados pelo projeto.

A camada Gold implementa o modelo relacional LegisAnalytica v5, com 26 tabelas normalizadas no schema `workspace.pi_ii_gold`, alimentadas a partir das 13 tabelas da camada Silver (`workspace.pi_ii_silver`).

## Notebook

- **LegisAnalytica - Gold Layer DDL** — cria as 26 tabelas Gold via `CREATE OR REPLACE TABLE` a partir das tabelas Silver, com verificação de contagem de linhas ao final.

## Entrada
- Dados tratados da camada Silver (`workspace.pi_ii_silver`)
- Modelo relacional de referência: `Modelagem_Relacional_LegisAnalytica.pdf`

## Saída
- 26 tabelas relacionais no schema `workspace.pi_ii_gold`, organizadas em quatro grupos:

| Grupo | Tabelas |
|-------|---------|
| Catálogos | CASA, UF, LEGISLATURA, PARTIDO, LOTE, CAUSA_AFASTAMENTO |
| Identidade & Exercício | PARLAMENTAR, IDENTIFICADOR_PARLAMENTAR, CADASTRO_OBSERVADO, PARLAMENTAR_LEGISLATURA, MANDATO, VINCULO_MANDATO, EXERCICIO |
| Atividade Legislativa | TIPO_PROPOSICAO, PROPOSICAO, AUTORIA, STATUS_ATUAL, VOTACAO, VOTACAO_PROPOSICAO, TIPO_VOTO, VOTO, TOTAL_PUBLICADO |
| Despesas | TIPO_DESPESA, FORNECEDOR, DESPESA, VALOR_DESPESA |

## Convenções de modelagem

- IDs são `STRING` (domínio do projeto)
- `casa` ∈ {CAMARA, SENADO}
- Zero em numero/ano → convertido em `NULL` via `NULLIF`
- `id_lote` = nome da tabela Silver de origem (chave para a tabela LOTE de proveniência)
- Booleans: Câmara `aprovacao = 1 → TRUE`; Senado `resultado LIKE '%aprovad%' → TRUE`
- Filtro `registro_valido_leg57 = true` aplicado nas tabelas do Senado quando aplicável

## Correções de integridade (27/10/2026)

3 correções pontuais aplicadas no notebook Gold DDL:

1. **`proposicao` (SENADO)**: adicionado `SELECT DISTINCT` para eliminar 2.060 chaves PK duplicadas — a Silver `materias` tem uma linha por matéria+autor (`dropDuplicates(["id_materia","id_senador"])`), e o Gold seleciona apenas colunas da matéria.
2. **`votacao_proposicao`**: adicionado `INNER JOIN` com a Gold `proposicao` para remover 980 registros órfãos (491 `id_proposicao=0` na Câmara + 73 proposições ausentes + 416 matérias inexistentes/fora da 57ª legislatura no Senado).
3. **`total_publicado` (SENADO)**: trocado `COUNT(*)` por `COUNT(DISTINCT id_materia)` para refletir proposições únicas.

**Impacto:** `proposicao` 264.078→248.560; `votacao_proposicao` 2.020→1.040; `total_publicado` SENADO 29.193→13.675.

## Limitações conhecidas (para PostgreSQL/Supabase)

- `despesa.codigo_tipo`: 329 registros do Senado com `NULL` (API do Senado não envia o tipo). FK `despesa→tipo_despesa` deve ser **nullable**.
- `despesa.id_parlamentar` (Câmara): é código **CEAP**, não FK direta para `parlamentar`. Mapeamento via `identificador_parlamentar` (`sistema='CEAP'`). 216 CEAPs sem correspondência → `id_parlamentar=NULL`.
- Todas as colunas são `nullable` no schema Gold (sem constraints PK/FK a nível DB). Constraints serão criadas no PostgreSQL.
- `votacao_proposicao` (SENADO): apenas 7 registros (votações que referenciam matérias válidas da 57ª legislatura).

## Migração para Supabase (PostgreSQL)

A camada Gold pode ser exportada para o PostgreSQL do Supabase para consumo pela aplicação web.

| Diretório | Arquivo | Descrição |
|-----------|--------|------------|
| `export/` | `EXPORT_supabase_postgresql.ipynb` | Notebook de transferência em lotes + validação |
| `supabase/migrations/` | `001_initial_schema_supabase.sql` | Esquema PostgreSQL: 26 tabelas, 60 FKs, 30 índices, 2 UNIQUEs |

**Fluxo de migração:**
- O esquema é criado no Supabase executando `001_initial_schema_supabase.sql` diretamente no SQL Editor (aplicação manual; o script é versionado aqui como referência). O notebook de exportação também verifica a existência do schema e o aplica automaticamente caso o banco esteja vazio.
- O notebook `EXPORT_supabase_postgresql.ipynb` opera em três modos, controlados pela variável `EXPORT_MODE`: `prepare` (padrão, sem gravação), `pilot` (apenas `public.casa`) e `full` (as 26 tabelas). O recorte de exportação é **2025 completo** (01/01/2025 a 31/12/2025), preservando entidades anteriores necessárias às FKs.
- Antes da primeira escrita, o notebook valida as constraints da migration (NOT NULL, PKs, UNIQUE, VARCHAR) e o schema de destino; após a carga, compara contagens por tabela e executa 46 verificações de integridade referencial.

**Pré-requisitos:** Secret Scope `legisanalytica_supabase` no Databricks com chaves `jdbc_url`, `user`, `password`.

## Observação
Alterações podem ocorrer ao longo do projeto. Atualizar o readme conforme necessário.
