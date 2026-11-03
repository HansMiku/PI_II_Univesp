package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"legisanalytica/internal/activity"
	"legisanalytica/internal/expense"
	"legisanalytica/internal/graphql"
	"legisanalytica/internal/politician"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestServer(t *testing.T) *httptest.Server {
	envPath := filepath.Join("../../", ".env")
	if err := godotenv.Load(envPath); err != nil {
		t.Logf("Warning: Could not load .env from %s: %v", envPath, err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL is required for integration tests")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	pRepo := politician.NewRepository(db)
	pSvc := politician.NewService(pRepo)
	eRepo := expense.NewRepository(db)
	eSvc := expense.NewService(eRepo)
	aSvcRepo := activity.NewRepository(db)
	aSvc := activity.NewService(aSvcRepo)

	handler := graphql.NewGraphQLHandler(pSvc, eSvc, aSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", handler.HandleGraphQL)

	return httptest.NewServer(mux)
}

func executeQuery(ts *httptest.Server, query string) (map[string]interface{}, error) {
	reqBody, _ := json.Marshal(map[string]string{"query": query})
	resp, err := http.Post(ts.URL+"/graphql", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func TestGraphQLIntegration(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	listQuery := "{ parlamentares { id_parlamentar casa } }"
	listRes, err := executeQuery(ts, listQuery)
	if err != nil {
		t.Fatalf("Failed to execute list query: %v", err)
	}

	data, ok := listRes["data"].(map[string]interface{})
	if !ok || data["parlamentares"] == nil {
		t.Skip("No politicians found in database, skipping detailed tests")
		return
	}

	politicians := data["parlamentares"].([]interface{})
	firstP := politicians[0].(map[string]interface{})
	id := firstP["id_parlamentar"].(string)
	casa := firstP["casa"].(string)

	tests := []struct {
		name  string
		query string
	}{
		{"ListAll", "{ parlamentares { nome_referencia } }"},
		{"GetOne", fmt.Sprintf("{ parlamentar(id: \"%s\", casa: \"%s\") { nome_referencia } }", id, casa)},
		{"GetExpenses", fmt.Sprintf("{ despesas(id: \"%s\") { valor_liquido } }", id)},
		{"GetExpensesWithDetails", fmt.Sprintf("{ despesas(id: \"%s\") { valor_liquido tipo { descricao } fornecedor { nome } } }", id)},
		{"GetActivity", fmt.Sprintf("{ atividade(id: \"%s\") { ementa } }", id)},
		{"GetVotes", fmt.Sprintf("{ votos(id: \"%s\") { id_votacao tipo { descricao } } }", id)},
		{"FullQuery", fmt.Sprintf("{ parlamentar(id: \"%s\", casa: \"%s\") { nome_referencia } despesas(id: \"%s\") { valor_liquido } }", id, casa, id)},
		{"InvalidID", "{ parlamentar(id: \"NON_EXISTENT\", casa: \"NON_EXISTENT\") { nome_referencia } }"},
		{"MissingArg", "{ parlamentar(id: \"some_id\") { nome_referencia } }"},
		{"OnlyEmails", "{ parlamentares { email } }"},
		{"ExpenseDetails", fmt.Sprintf("{ despesas(id: \"%s\") { cnpj_cpf } }", id)},
		{"ActivityDetails", fmt.Sprintf("{ atividade(id: \"%s\") { sigla_tipo } }", id)},
		{"SupportTables", "{ casas { nome } ufs { nome } partidos { sigla } legislaturas { numero_legislatura } }"},
		{"GlobalPropositions", "{ todas_proposicoes(limit: 5) { ementa } }"},
		{"GlobalVotes", "{ todos_votos(limit: 5) { id_votacao } }"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := executeQuery(ts, tt.query)
			if err != nil {
				t.Errorf("Request failed: %v", err)
			}
			if res == nil {
				t.Error("Response is nil")
			}
		})
	}
}
