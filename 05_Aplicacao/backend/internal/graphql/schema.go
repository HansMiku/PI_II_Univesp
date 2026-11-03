package graphql

import (
	"github.com/graphql-go/graphql"

	"legisanalytica/internal/activity"
	"legisanalytica/internal/expense"
	"legisanalytica/internal/politician"
)

type schemaBuilder struct {
	ps politician.Service
	ds expense.Service
	as activity.Service

	houseType       *graphql.Object
	stateType       *graphql.Object
	partyType       *graphql.Object
	legislatureType *graphql.Object

	expenseKindType *graphql.Object
	supplierType    *graphql.Object
	expenseType     *graphql.Object

	propositionType *graphql.Object
	votacaoType     *graphql.Object
	voteKindType    *graphql.Object
	voteType        *graphql.Object

	politicianType *graphql.Object
}

func (b *schemaBuilder) rootQuery() *graphql.Object {
	fields := graphql.Fields{}
	for _, group := range []graphql.Fields{
		b.referenceQueries(),
		b.politicianQueries(),
		b.expenseQueries(),
		b.activityQueries(),
	} {
		for name, f := range group {
			fields[name] = f
		}
	}
	return graphql.NewObject(graphql.ObjectConfig{Name: "RootQuery", Fields: fields})
}

func buildSchema(ps politician.Service, ds expense.Service, as activity.Service) (graphql.Schema, error) {
	b := &schemaBuilder{ps: ps, ds: ds, as: as}

	b.buildReferenceTypes()
	b.buildExpenseTypes()
	b.buildActivityTypes()
	b.buildPoliticianType()

	return graphql.NewSchema(graphql.SchemaConfig{Query: b.rootQuery()})
}
