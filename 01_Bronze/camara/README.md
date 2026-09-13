# Bronze Câmara - Brasil por UF (v3)

Pipeline:
00 -> (01, 02 e 03 em paralelo) -> 04

Ajustes desta versão:
- contagem de deputados separa registros e IDs únicos;
- fallback de proposições somente para IDs sem detalhe no CSV anual;
- arquivos `proposicoes_fallback_YYYY.csv` por UF;
- auditoria de despesas não distribuídas em `_despesas_nao_distribuidas_YYYY.csv`;
- validação atualizada;
- ano corrente marcado como parcial;
- diretórios temporários exclusivos por execução, mantendo compatibilidade com Serverless/retry.

Destino:
`/Volumes/workspace/pi_ii_bronze/camara/<uf>/`
