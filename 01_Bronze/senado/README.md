# Bronze Senado - Brasil por UF (v2)

Correção da população histórica da 57ª legislatura.

## Job

```text
                         ┌─ 01_materias ─┐
00_setup_senadores ──────┼─ 02_despesas ┼── 04_validacao
                         └─ 03_votacoes ─┘
```

## O que mudou

- o notebook 00 usa a lista da 57ª legislatura, mas só mantém exercícios que
  realmente intersectam 01/02/2023 a 31/01/2027;
- histórico de mandatos antigos carregado no retorno da API não entra
  automaticamente no universo;
- suplentes que efetivamente exerceram entram normalmente;
- `senadores_exercicios.csv` registra os períodos de exercício por UF;
- autorias são consultadas apenas para a coorte histórica válida e filtradas
  pelo período de exercício;
- CEAPS usa o mapa histórico corrigido;
- despesas sem correspondência não são apagadas silenciosamente: ficam nos
  arquivos `_despesas_nao_distribuidas_YYYY.csv` e fazem a validação retornar
  `revisar`;
- votações continuam usando a UF do próprio registro do XML anual.

Destino:

`/Volumes/workspace/pi_ii_bronze/senado/<uf>/`


## v2.1

Corrige o notebook 03. A versão anterior ainda chamava uma função da v1 que procurava
as colunas antigas em `senadores.csv`. Agora o fallback de UF usa diretamente
`_senadores_exercicios_brasil.csv`.


## v2.2

Corrige a leitura de datas do Senado.

- datas no formato `YYYY-MM-DD` agora são tratadas explicitamente como ISO;
- mantém fallback `dayfirst=True` apenas para formatos não ISO;
- valida e interrompe o setup se algum período ficar com início depois do fim;
- aplica o mesmo parser em `Materia_Data`;
- a validação final também confere períodos históricos invertidos.
