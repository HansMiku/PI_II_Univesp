package graphql

import "github.com/graphql-go/graphql"

func (b *schemaBuilder) referenceQueries() graphql.Fields {
	return graphql.Fields{
		"casas": {
			Type: graphql.NewList(b.houseType),

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return b.ps.ListAllHouses(p.Context)
			},
		},

		"ufs": {
			Type: graphql.NewList(b.stateType),

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return b.ps.ListAllStates(p.Context)
			},
		},

		"partidos": {
			Type: graphql.NewList(b.partyType),

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return b.ps.ListAllParties(p.Context)
			},
		},

		"legislaturas": {
			Type: graphql.NewList(b.legislatureType),

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return b.ps.ListAllLegislatures(p.Context)
			},
		},
	}
}
func (b *schemaBuilder) buildReferenceTypes() {
	b.houseType = graphql.NewObject(graphql.ObjectConfig{
		Name: "House",
		Fields: graphql.Fields{
			"casa": &graphql.Field{Type: graphql.String},
			"nome": &graphql.Field{Type: graphql.String},
		},
	})

	b.stateType = graphql.NewObject(graphql.ObjectConfig{
		Name: "State",
		Fields: graphql.Fields{
			"uf":   &graphql.Field{Type: graphql.String},
			"nome": &graphql.Field{Type: graphql.String},
		},
	})

	b.partyType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Party",
		Fields: graphql.Fields{
			"sigla": &graphql.Field{Type: graphql.String},
		},
	})

	b.legislatureType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Legislature",
		Fields: graphql.Fields{
			"numero_legislatura": &graphql.Field{Type: graphql.Int},
			"inicio":             &graphql.Field{Type: graphql.String},
			"fim":                &graphql.Field{Type: graphql.String},
		},
	})
}
