package api

import (
	"strings"

	"github.com/graphql-go/graphql"
)

func (b *schemaBuilder) politicianQueries() graphql.Fields {
	return graphql.Fields{
		"parlamentar": &graphql.Field{
			Type: b.politicianType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"casa": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(string)
				house := p.Args["casa"].(string)

				return b.ps.GetPoliticianFullData(
					p.Context,
					id,
					house,
				)
			},
		},
		"parlamentares": &graphql.Field{
			Type: graphql.NewList(b.politicianType),
			Args: graphql.FieldConfigArgument{
				"limit":  &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
				"offset": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 0},
				"uf":     &graphql.ArgumentConfig{Type: graphql.String},
				"casa":   &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				limit, offset := paginationArgs(p)

				uf, _ := p.Args["uf"].(string)
				if uf == "" {
					return b.ps.ListAllPoliticians(p.Context, limit, offset)
				}

				var house *string
				if c, ok := p.Args["casa"].(string); ok && c != "" {
					house = &c
				}
				return b.ps.ListPoliticiansByState(p.Context, strings.ToUpper(uf), house, limit, offset)
			},
		},
		"buscaParlamentares": &graphql.Field{
			Type: graphql.NewList(b.politicianType),
			Args: graphql.FieldConfigArgument{
				"nome":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"limit":  &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
				"offset": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 0},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name := p.Args["nome"].(string)
				limit, offset := paginationArgs(p)

				return b.ps.SearchPoliticiansByName(p.Context, name, limit, offset)
			},
		},
	}
}

func (b *schemaBuilder) buildPoliticianType() {
	b.politicianType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Politician",
		Fields: graphql.Fields{
			"casa":            &graphql.Field{Type: graphql.String},
			"id_parlamentar":  &graphql.Field{Type: graphql.String},
			"nome_referencia": &graphql.Field{Type: graphql.String},
			"url_foto":        &graphql.Field{Type: graphql.String},
			"email":           &graphql.Field{Type: graphql.String},

			"estado": &graphql.Field{
				Type: b.stateType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					pol := politicianFromSource(p.Source)
					if pol == nil {
						return nil, nil
					}
					reg, err := b.ps.GetObservedRegistration(p.Context, pol.House, pol.IDPolitician)
					if err != nil || reg == nil {
						return nil, nil
					}
					return b.ps.GetState(p.Context, reg.UF)
				},
			},

			"partido": &graphql.Field{
				Type: b.partyType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					pol := politicianFromSource(p.Source)
					if pol == nil {
						return nil, nil
					}
					reg, err := b.ps.GetObservedRegistration(p.Context, pol.House, pol.IDPolitician)
					if err != nil || reg == nil {
						return nil, nil
					}
					return b.ps.GetParty(p.Context, reg.Party)
				},
			},

			"legislaturas": &graphql.Field{
				Type: graphql.NewList(b.legislatureType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					pol := politicianFromSource(p.Source)
					if pol == nil {
						return nil, nil
					}
					pls, err := b.ps.GetPoliticianLegislatures(p.Context, pol.House, pol.IDPolitician)
					if err != nil {
						return nil, err
					}
					var results []interface{}
					for _, pl := range pls {
						leg, err := b.ps.GetLegislature(p.Context, pl.NumberLegislature)
						if err == nil {
							results = append(results, leg)
						}
					}
					return results, nil
				},
			},

			"despesas": &graphql.Field{
				Type: graphql.NewList(b.expenseType),
				Args: graphql.FieldConfigArgument{
					"limit":     limitArg(),
					"offset":    offsetArg(),
					"type_code": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					pol := politicianFromSource(p.Source)
					if pol == nil {
						return nil, nil
					}
					limit, offset := paginationArgs(p)
					var typeCode *string
					if val, ok := p.Args["type_code"].(string); ok {
						typeCode = &val
					}
					return b.ds.GetExpensesByPolitician(p.Context, pol.IDPolitician, limit, offset, typeCode)
				},
			},

			"atividade": &graphql.Field{
				Type: graphql.NewList(b.propositionType),
				Args: graphql.FieldConfigArgument{
					"limit":  limitArg(),
					"offset": offsetArg(),
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					pol := politicianFromSource(p.Source)
					if pol == nil {
						return nil, nil
					}
					limit, offset := paginationArgs(p)
					return b.as.GetPropositionsByPolitician(p.Context, pol.IDPolitician, limit, offset)
				},
			},

			"votos": &graphql.Field{
				Type: graphql.NewList(b.voteType),
				Args: graphql.FieldConfigArgument{
					"limit":     limitArg(),
					"offset":    offsetArg(),
					"vote_code": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					pol := politicianFromSource(p.Source)
					if pol == nil {
						return nil, nil
					}
					limit, offset := paginationArgs(p)
					var voteCode *string
					if val, ok := p.Args["vote_code"].(string); ok {
						voteCode = &val
					}
					return b.as.GetVotesByPolitician(p.Context, pol.IDPolitician, limit, offset, voteCode)
				},
			},
		},
	})
}
