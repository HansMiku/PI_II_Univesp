# Analytics — Análise de Dados

Responsável: Vitória

Esta pasta contém os códigos utilizados para geração de indicadores e análises a partir dos dados consolidados na camada Gold.

## Entrada
- Dados da camada Gold

## Saída
- Indicadores e resultados analíticos utilizados pela aplicação e pelas visualizações

## Observação
Alterações podem ocorrer ao longo do projeto. Atualizar o readme conforme necessário.

# Sugestões para Analytics

A etapa de Analytics deve consumir os indicadores e agregações produzidos na camada Gold.

O objetivo é consultar, validar e explorar os indicadores definidos no projeto, preparando os dados que poderão ser disponibilizados pelo backend e posteriormente apresentados no frontend.

## Atividade parlamentar

- consultar o número de proposições apresentadas;
- comparar proposições por tipo e período;
- consultar o número de votações participadas;
- analisar taxa de participação e ausências, caso esses indicadores sejam considerados metodologicamente válidos.

## Comportamento em votações

- analisar percentuais de votos "Sim", "Não" e abstenções;
- comparar a distribuição dos votos entre partidos;
- comparar comportamento de voto entre parlamentares, partidos e UFs;
- analisar alinhamento partidário, caso o indicador seja implementado na Gold.

## Produção legislativa

- comparar o total de proposições entre parlamentares;
- comparar produção por partido e UF;
- analisar a evolução temporal da produção;
- consultar situação das proposições;
- analisar aprovação ou avanço das proposições, caso o indicador seja implementado.

## Despesas

- comparar despesas entre parlamentares;
- comparar despesas entre partidos e UFs;
- analisar despesas por categoria;
- analisar a evolução mensal dos gastos;
- identificar valores muito acima ou abaixo do padrão observado.

## Representação política

- analisar a distribuição dos parlamentares por partido;
- analisar a representação por UF;
- comparar a composição partidária das Casas;
- comparar Câmara e Senado.

## Indicadores estatísticos

- analisar média, mediana, mínimo e máximo;
- analisar desvio-padrão e percentis;
- analisar proporções e variações temporais;
- gerar rankings descritivos;
- investigar correlações entre variáveis quando houver sentido analítico.

## Saídas da etapa

A etapa pode produzir:

- consultas SQL/DQL;
- views ou consultas preparadas para consumo;
- tabelas de validação;
- análises exploratórias;
- documentação das regras de consulta;
- protótipos de visualização para validação dos indicadores, quando necessário.

As visualizações finais da aplicação devem ser implementadas no frontend, utilizando os dados disponibilizados pelo backend.

Os cálculos principais e as regras dos indicadores devem permanecer na camada Gold, evitando duplicação da lógica de negócio.