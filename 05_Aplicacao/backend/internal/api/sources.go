package api

import (
	"legisanalytica/internal/activity"
	"legisanalytica/internal/expense"
	"legisanalytica/internal/politician"
)

func politicianFromSource(src interface{}) *politician.Politician {
	switch s := src.(type) {
	case *politician.Politician:
		return s
	case politician.Politician:
		return &s
	}
	return nil
}

func expenseFromSource(src interface{}) *expense.Expense {
	switch s := src.(type) {
	case *expense.Expense:
		return s
	case expense.Expense:
		return &s
	}
	return nil
}

func votacaoFromSource(src interface{}) *activity.Votacao {
	switch s := src.(type) {
	case *activity.Votacao:
		return s
	case activity.Votacao:
		return &s
	}
	return nil
}

func voteFromSource(src interface{}) *activity.Vote {
	switch s := src.(type) {
	case *activity.Vote:
		return s
	case activity.Vote:
		return &s
	}
	return nil
}
