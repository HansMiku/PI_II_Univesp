package api

import "github.com/graphql-go/graphql"

func (b *schemaBuilder) expenseQueries() graphql.Fields {
	return graphql.Fields{
		"despesas": &graphql.Field{
			Type: graphql.NewList(b.expenseType),

			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"limit": {
					Type:         graphql.Int,
					DefaultValue: 20,
				},
				"offset": {
					Type:         graphql.Int,
					DefaultValue: 0,
				},
				"type_code": {
					Type: graphql.String,
				},
			},

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(string)

				limit, offset := paginationArgs(p)

				var typeCode *string
				if val, ok := p.Args["type_code"].(string); ok {
					typeCode = &val
				}

				return b.ds.GetExpensesByPolitician(
					p.Context,
					id,
					limit,
					offset,
					typeCode,
				)
			},
		},
	}
}
func (b *schemaBuilder) buildExpenseTypes() {
	b.expenseKindType = graphql.NewObject(graphql.ObjectConfig{
		Name: "ExpenseType",
		Fields: graphql.Fields{
			"casa":        &graphql.Field{Type: graphql.String},
			"codigo_tipo": &graphql.Field{Type: graphql.String},
			"descricao":   &graphql.Field{Type: graphql.String},
		},
	})

	b.supplierType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Supplier",
		Fields: graphql.Fields{
			"casa":     &graphql.Field{Type: graphql.String},
			"cnpj_cpf": &graphql.Field{Type: graphql.String},
			"nome":     &graphql.Field{Type: graphql.String},
		},
	})

	b.expenseType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Expense",
		Fields: graphql.Fields{
			"casa":            &graphql.Field{Type: graphql.String},
			"id_despesa":      &graphql.Field{Type: graphql.String},
			"id_parlamentar":  &graphql.Field{Type: graphql.String},
			"ano":             &graphql.Field{Type: graphql.Int},
			"mes":             &graphql.Field{Type: graphql.Int},
			"data_emissao":    &graphql.Field{Type: graphql.String},
			"codigo_tipo":     &graphql.Field{Type: graphql.String},
			"cnpj_cpf":        &graphql.Field{Type: graphql.String},
			"valor_documento": &graphql.Field{Type: graphql.Float},
			"valor_glosa":     &graphql.Field{Type: graphql.Float},
			"valor_liquido":   &graphql.Field{Type: graphql.Float},
			"tipo": &graphql.Field{
				Type: b.expenseKindType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					exp := expenseFromSource(p.Source)
					if exp == nil {
						return nil, nil
					}
					return b.ds.GetExpenseType(p.Context, exp.House, exp.TypeCode)
				},
			},
			"fornecedor": &graphql.Field{
				Type: b.supplierType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					exp := expenseFromSource(p.Source)
					if exp == nil {
						return nil, nil
					}
					return b.ds.GetSupplier(p.Context, exp.House, exp.TaxID)
				},
			},
		},
	})
}
