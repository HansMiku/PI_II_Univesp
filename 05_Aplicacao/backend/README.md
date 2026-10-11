# Backend

Responsável: João

Esta pasta contém o backend da aplicação e os endpoints utilizados para disponibilizar os dados e indicadores ao frontend.

Credenciais, tokens e outras informações sensíveis não devem ser versionados no GitHub.

## Fonte de dados

O backend deve consumir as 26 tabelas da camada Gold no schema `workspace.pi_ii_gold`.

### Convenções importantes

- Todos os IDs são `STRING` (não INT) — considerar na serialização JSON e nos tipos de resposta da API
- `casa` ∈ {`'CAMARA'`, `'SENADO'`} — pode ser usado como parâmetro de filtro nos endpoints
- Chaves primárias compostas são comuns (ex: `casa + id_parlamentar`, `casa + id_proposicao`)
- Para despesas, a tabela `despesa` traz o cabeçalho e `valor_despesa` traz os valores monetários; o JOIN é por `casa + id_despesa`
- Para votações, a tabela `votacao_proposicao` faz a ligação entre `votacao` e `proposicao`

### Tabelas mais relevantes para a API

| Endpoint provável | Tabelas Gold |
|------------------|--------------|
| Parlamentares | `parlamentar`, `cadastro_observado`, `parlamentar_legislatura` |
| Proposições | `proposicao`, `autoria`, `tipo_proposicao`, `status_atual` |
| Votações | `votacao`, `voto`, `tipo_voto`, `votacao_proposicao` |
| Despesas | `despesa`, `valor_despesa`, `tipo_despesa`, `fornecedor` |
| Filtros | `casa`, `uf`, `partido`, `legislatura` |

## Observação
Alterações podem ocorrer ao longo do projeto. Atualizar o readme conforme necessário.