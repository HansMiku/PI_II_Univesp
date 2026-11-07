# Frontend — LegisData

Responsáveis:

- Brenda — interface e desenvolvimento frontend
- Eder — visualização e acessibilidade

Esta pasta contém a interface da aplicação, incluindo páginas, componentes visuais, gráficos e recursos de acessibilidade.

## Páginas

| Arquivo | Conteúdo |
|---|---|
| `index.html` | Página inicial, com apresentação do projeto e links para os painéis |
| `despesas.html` + `despesas.js` | Painel de despesas por parlamentar, com resumo por categoria e filtro de período |
| `propostas.html` + `propostas.js` | Painel de propostas de autoria dos parlamentares, com situação de votação |
| `style.css` | Identidade visual, paleta, componentes e responsividade |

## Como executar

### 1. Backend

```powershell
cd 05_Aplicacao\backend
go run ./cmd/api
```

A API sobe em `http://127.0.0.1:8080/graphql`.

### 2. Frontend

Abra a pasta `05_Aplicacao/frontend` com um servidor estático. O projeto foi validado com o Live Server em:

```text
http://127.0.0.1:5501/05_Aplicacao/frontend/
```

A URL da API está definida diretamente em `despesas.js` e `propostas.js`:

```js
const API_URL = 'http://127.0.0.1:8080/graphql';
```

### 3. CORS

O backend precisa permitir a origem usada pelo frontend. Exemplo de `.env`:

```env
CORS_ALLOWED_ORIGIN=http://127.0.0.1:5501
```

Se a porta do Live Server mudar, ajuste também o `CORS_ALLOWED_ORIGIN` e reinicie a API.


## Identidade visual

O nome adotado em todas as páginas é **LegisData**.

A paleta está centralizada em variáveis no `style.css`:

| Variável | Uso |
|---|---|
| `--cor-marinho` | Navbar, títulos e base da identidade |
| `--cor-ciano` | Destaques, links ativos e detalhes |
| `--cor-verde` | Botões primários e valores positivos |
| `--cor-texto` | Texto principal |
| `--cor-texto-secundario` | Textos de apoio |
| `--cor-borda` | Bordas e divisores |
| `--cor-aviso-*` | Avisos sobre qualidade dos dados |

As páginas internas (`despesas.html` e `propostas.html`) usam o mesmo padrão visual da home:

- seção de destaque com degradê suave e formas geométricas;
- selo de identificação do painel;
- título e descrição em alto contraste;
- pílula de status da API;
- cards com borda, sombra e cantos arredondados.

## Integração com a API

### Despesas

Consultas usadas pela página:

| Consulta | Finalidade |
|---|---|
| `parlamentares(limit, offset)` | Popular a busca de parlamentares |
| `resumoDespesas(id, casa, mes_inicio, mes_fim)` | Resumo agregado por categoria |
| `estatisticasDespesas(id, casa, mes_inicio, mes_fim)` | Total, quantidade, meses com despesas e média mensal |

Consultas disponíveis, mas não usadas no resumo atual:

| Consulta | Finalidade |
|---|---|
| `despesas(id, casa, limit, offset, type_code)` | Listagem paginada de despesas individuais |
| `totalDespesas(id, casa, type_code)` | Contagem total de despesas |

Observações:

- Na Câmara, `despesa.id_parlamentar` é o código CEAP. O backend traduz para o id da API usando `identificador_parlamentar`.
- O resumo usa **todas as despesas** do parlamentar no período, sem paginação.
- O gráfico pode exibir **total do período** ou **média mensal** (total dividido pelos meses com despesas).
- Se não houver despesas, a página mostra uma mensagem e oculta o gráfico.

### Propostas

Consultas usadas pela página:

| Consulta | Finalidade |
|---|---|
| `parlamentares(limit, offset)` | Popular a busca de parlamentares |
| `atividade(id, limit, offset, somente_votadas)` | Listar propostas de autoria do parlamentar |
| `votacoes_da_proposicao(casa, id_proposicao)` | Consultar a situação e os votos da proposta |

Observações:

- A lista mostra apenas propostas em que o parlamentar é **autor principal**.
- O campo `tem_votacao` indica se existe votação registrada.
- O filtro **“Mostrar somente propostas com votação registrada”** usa `somente_votadas: true`.
- Propostas com votação recebem o selo **Votada**.
- O modal mostra:
  - autor principal;
  - situação da votação (`Aprovada na votação` / `Rejeitada na votação`);
  - data, órgão, descrição e contagem de votos;
  - mensagem específica quando não há votação registrada.
- `numero` e `ano` podem vir `null`; a interface trata esses casos sem exibir `0`.

## Alterações de backend que impactam o frontend

### Despesas

- `despesas` passou a aceitar o argumento opcional `casa`.
- O join com `valor_despesa` passou a considerar `casa + id_despesa`, evitando duplicidades.
- A Câmara passou a usar o mapeamento CEAP → id da API.
- Foram adicionadas as consultas `resumoDespesas`, `estatisticasDespesas` e `totalDespesas`.
- A listagem ganhou ordenação estável para paginação.

### Propostas

- Foi adicionada a coluna `autoria.principal`:
  - Câmara: primeiro signatário (`ordem_assinatura = 1`);
  - Senado: `IndicadorAutorPrincipal = "Sim"`.
- A consulta `atividade` filtra `autoria.principal = true`.
- A consulta `atividade` ganhou `somente_votadas`.
- O tipo `Proposition` ganhou `tem_votacao`.
- `numero` e `ano` agora podem ser nulos, em vez de `0`.
- A consulta `votacoes_da_proposicao` é usada para exibir a situação da proposta.

### Dados e exportação

- O notebook `03_Gold/LegisAnalytica - Gold Layer DDL.ipynb` grava a coluna `principal` em `autoria`.
- A migration `03_Gold/supabase/migrations/002_add_autoria_principal.sql` adiciona a coluna na base PostgreSQL.
- O comando de backfill da base existente é:

```powershell
cd 05_Aplicacao\backend
go run ./cmd/backfill-autoria
```

- O notebook `03_Gold/export/EXPORT_supabase_postgresql.ipynb` foi ajustado para exportar:
  - votações de 2025;
  - votações vinculadas às proposições exportadas;
  - votos dessas votações.
- Para refletir a ampliação de cobertura no Supabase, é necessário executar novamente a exportação Gold → Supabase.

## Limitações conhecidas

- A base atual de despesas cobre o recorte de **2025**.
- Alguns parlamentares podem não ter mapeamento CEAP na Câmara; nesses casos, não há despesas vinculadas.
- A quantidade de propostas com votação registrada ainda é limitada na base atual. O selo e o filtro **Votada** ajudam a localizar as que existem.
- O Senado pode ter poucos ou nenhum vínculo em `votacao_proposicao` na base atual; a exportação foi ajustada para recuperar os vínculos existentes no Gold.
- Propostas sem número ou ano são exibidas sem o valor `0`.
- A aprovação só é exibida quando existe uma votação formal vinculada. Aprovações por despacho ou etapas sem votação registrada não aparecem.

## Checklist antes de rodar

1. API em `http://127.0.0.1:8080/graphql`.
2. Frontend em `http://127.0.0.1:5501`.
3. `CORS_ALLOWED_ORIGIN` igual à origem do frontend.
4. `.env` fora do Git.
5. Se a API estiver em outro endereço, ajustar `API_URL` nos dois arquivos JS.


