# Export — Gold → Supabase

Notebook de transferência de `workspace.pi_ii_gold` para o PostgreSQL Supabase.

## Arquivos

- `EXPORT_supabase_postgresql.ipynb`: notebook para Databricks (versão a publicar no GitHub).

## Pré-requisitos

1. Executar `../supabase/migrations/001_initial_schema_supabase.sql` no SQL Editor do projeto de destino (apenas para banco vazio). O notebook também verifica se o schema existe e aplica a migration automaticamente quando o banco está vazio.
2. Secret Scope `legisanalytica_supabase` com `jdbc_url`, `user` e `password` atualizados.
3. Driver `psycopg` v3 no Databricks, instalado separadamente, com `PSYCOPG_IMPL=python` antes do import.
4. Permissão de leitura na Gold e conexão funcional com o novo projeto Supabase.

## Recorte e relacionamentos

- **01/01/2025 a 31/12/2025** (um ano-calendário completo).
- Despesas e total publicado por `ano=2025`; votações por `data_votacao`.
- Proposições apresentadas em 2025 **mais** proposições anteriores referenciadas por votações de 2025, necessárias para preservar FKs. Algumas proposições podem ter `data_apresentacao` anterior a 2025: não classificá-las como novas proposições do ano.
- `autoria`, `status_atual` e `votacao_proposicao` seguem as entidades selecionadas; `valor_despesa` segue as despesas e `voto` segue as votações.
- Catálogos e relacionamentos cadastrais permanecem completos quando necessários.
- Votações sem `data_votacao` não entram no recorte por falta de data para classificação.

## Procedimento seguro

1. O notebook inicia com `EXPORT_MODE = "prepare"` (**não escreve** no banco).
2. Para testar a inserção de dois registros de `casa`, trocar para `EXPORT_MODE = "pilot"`.
3. Conferir `public.casa`, o tamanho do PostgreSQL e a métrica **Disk** no painel do Supabase.
4. Somente depois, trocar para `EXPORT_MODE = "full"` para exportar as 26 tabelas.
5. Antes da primeira escrita, o notebook valida as constraints da migration (NOT NULL, PKs, UNIQUE, VARCHAR) e o schema de destino; qualquer violação aborta a execução sem gravar.
6. O notebook utiliza lotes de 5.000 e `ON CONFLICT DO NOTHING`: evita duplicações em novas tentativas, **mas não atualiza registros já existentes**.
7. Falhas de lote interrompem a carga; commits anteriores são mantidos. Na execução integral, a validação compara contagens por tabela e executa 46 verificações de integridade referencial (FKs).

## Limites e cautelas

- O monitoramento consulta tamanho do banco e WAL, com limites conservadores e interrupção em caso de falha na medição. **Isso não mede todo o disco físico**. Consulte regularmente o painel Disk do Supabase; o próprio banco pode sofrer restrições de escrita.
- Não há garantia de que o recorte caiba no Free antes da medição após a carga-piloto.
- Não publicar o notebook temporário de cadastro de secrets, senhas, tokens, logs de conexão ou dados sensíveis no GitHub.
