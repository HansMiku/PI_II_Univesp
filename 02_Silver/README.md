# Silver — Qualidade e Transformação

Responsável: Rafael

Esta pasta contém os códigos de limpeza, padronização e transformação dos dados provenientes da camada Bronze.

## Entrada
- Dados da camada Bronze

## Saída
- Dados tratados e padronizados na camada Silver

## Observação
Alterações podem ocorrer ao longo do projeto. Atualizar o readme conforme necessário.

# Sugestões para a camada Silver

A camada Silver deve preparar e padronizar os dados coletados na Bronze para uso nas etapas seguintes.

Principais pontos sugeridos:

- padronizar os dados da Câmara e do Senado;
- definir identificadores consistentes para os parlamentares;
- padronizar Casa legislativa, UF, partido e legislatura;
- normalizar datas e tipos numéricos;
- padronizar os valores de voto;
- padronizar categorias de despesas;
- consolidar proposições/matérias e autorias;
- tratar registros duplicados e valores ausentes;
- preservar períodos de exercício dos parlamentares;
- manter rastreabilidade até os dados da Bronze;
- adicionar a legislatura como dimensão explícita.

A Silver deve concentrar os dados tratados e consistentes. Os cálculos dos indicadores devem ser realizados na camada Gold.