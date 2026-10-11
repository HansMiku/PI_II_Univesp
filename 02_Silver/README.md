# Silver — Qualidade e Transformação

Responsável: Rafael

Esta pasta contém os códigos de limpeza, padronização e transformação dos dados provenientes da camada Bronze.

## Notebook

- **Silver Layer Transformations** — lê os dados CSV da camada Bronze (volumes `workspace.pi_ii_bronze`), aplica padronização e tipos, e grava 13 tabelas no schema `workspace.pi_ii_silver`.

## Entrada
- Dados CSV da camada Bronze (`/Volumes/workspace/pi_ii_bronze/camara` e `/Volumes/workspace/pi_ii_bronze/senado`)

## Saída
- 13 tabelas tratadas e padronizadas no schema `workspace.pi_ii_silver`:

| Origem | Tabelas Silver |
|--------|----------------|
| Câmara | `deputados`, `proposicoes`, `proposicoes_autores`, `votacoes_camara`, `votos_camara`, `despesas_camara` |
| Senado | `senadores`, `senadores_exercicios`, `materias`, `materias_autorias`, `votacoes_senado`, `votos_senado`, `despesas_senado` |

Todas as tabelas incluem as colunas de partição `uf_particao` e `ano_particao` para rastreabilidade.

## Observação
Alterações podem ocorrer ao longo do projeto. Atualizar o readme conforme necessário.
