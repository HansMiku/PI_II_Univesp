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

## Observação
Alterações podem ocorrer ao longo do projeto. Atualizar o readme conforme necessário.
