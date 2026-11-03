package api

import "github.com/graphql-go/graphql"

const maxLimit = 100

func limitArg() *graphql.ArgumentConfig {
	return &graphql.ArgumentConfig{
		Type:        graphql.NewNonNull(graphql.Int),
		Description: "Obrigatório. Máximo 100.",
	}
}

func offsetArg() *graphql.ArgumentConfig {
	return &graphql.ArgumentConfig{
		Type:         graphql.Int,
		DefaultValue: 0,
	}
}

func paginationArgs(p graphql.ResolveParams) (limit, offset int) {
	limit, _ = p.Args["limit"].(int)
	offset, _ = p.Args["offset"].(int)

	if limit < 1 {
		limit = 1
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
