# Gold — Data Warehouse

Responsável: Bruce

Esta pasta contém os códigos relacionados à modelagem e organização dos dados consolidados utilizados pelo projeto.

A camada Gold será estruturada para facilitar consultas, análises, indicadores e consumo pela aplicação.

## Entrada
- Dados tratados da camada Silver

## Saída
- Tabelas, dimensões, fatos e views da camada Gold

## Observação
Alterações podem ocorrer ao longo do projeto. Atualizar o readme conforme necessário.

# Sugestões de indicadores para a camada Gold

A camada Gold deve utilizar os dados tratados na Silver para calcular indicadores e agregações que serão consumidos pela etapa de Analytics.

## Atividade parlamentar

- número de proposições apresentadas;
- proposições por tipo;
- proposições por período;
- taxa de participação em votações;
- número de votações participadas;
- ausências em votações.

## Comportamento em votações

- percentual de votos "Sim";
- percentual de votos "Não";
- percentual de abstenções;
- distribuição dos votos por partido;
- comparação do voto individual com a orientação/comportamento do partido;
- grau de alinhamento partidário.

## Produção legislativa

- total de proposições;
- proposições por parlamentar;
- proposições por partido;
- proposições por UF;
- evolução temporal da produção;
- situação das proposições;
- taxa de aprovação ou avanço das proposições, caso os dados permitam uma definição consistente.

## Despesas

- despesa total por parlamentar;
- despesa média;
- despesa por categoria;
- percentual de cada categoria no gasto total;
- evolução mensal;
- comparação entre parlamentares;
- comparação entre partidos;
- comparação entre UFs;
- distribuição e identificação de valores muito acima ou abaixo da média.

## Representação política

- número de parlamentares por partido;
- número de parlamentares por UF;
- participação percentual de cada partido na Casa;
- distribuição dos partidos por UF;
- comparação entre Câmara e Senado.

## Indicadores estatísticos

- média;
- mediana;
- mínimo e máximo;
- desvio-padrão;
- percentis;
- variação temporal;
- proporções;
- rankings descritivos;
- correlações entre variáveis, quando fizer sentido.

## Observações

Alguns indicadores exigem definição metodológica antes da implementação, principalmente:

- taxa de participação em votações;
- ausências em votações;
- orientação partidária;
- grau de alinhamento partidário;
- taxa de aprovação ou avanço das proposições.

Esses indicadores só devem ser calculados quando os dados disponíveis permitirem uma regra consistente e documentada.