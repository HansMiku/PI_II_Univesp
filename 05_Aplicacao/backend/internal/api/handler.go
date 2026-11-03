package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"legisanalytica/internal/activity"
	"legisanalytica/internal/expense"
	"legisanalytica/internal/politician"

	"github.com/graphql-go/graphql"
)

type GraphQLHandler struct {
	schema graphql.Schema
	logger *slog.Logger
}

func NewGraphQLHandler(ps politician.Service, ds expense.Service, as activity.Service) *GraphQLHandler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	schema, err := buildSchema(ps, ds, as)
	if err != nil {
		return nil
	}

	return &GraphQLHandler{schema: schema, logger: logger}
}

const (
	requestTimeout = 60 * time.Second
	maxBodyBytes   = 32 << 10
)

func (h *GraphQLHandler) HandleGraphQL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Query string `json:"query"`
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !validateQueryDepth(request.Query, 10) {
		http.Error(w, "Query depth exceeds maximum allowed limit", http.StatusBadRequest)
		return
	}

	if !validateQueryCost(request.Query) {
		http.Error(w, "Query complexity exceeds maximum allowed limit", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	result := graphql.Do(graphql.Params{
		Schema:        h.schema,
		RequestString: request.Query,
		Context:       ctx,
	})

	sanitizeErrors(h.logger, result.Errors)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
