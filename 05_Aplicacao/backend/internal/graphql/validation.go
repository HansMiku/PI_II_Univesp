package graphql

import (
	"strconv"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

const (
	maxRootFields = 10

	maxQueryCost = 5000

	maxCostWalkSteps = 100_000
	maxCostWalkDepth = 32
)

type queryComplexity struct {
	fragments  map[string]*ast.FragmentDefinition
	rootFields int
	cost       int
	steps      int
	aborted    bool
}

func validateQueryCost(query string) bool {
	doc, err := parser.Parse(parser.ParseParams{Source: query})
	if err != nil {
		return true
	}

	c := &queryComplexity{fragments: map[string]*ast.FragmentDefinition{}}
	for _, def := range doc.Definitions {
		if f, ok := def.(*ast.FragmentDefinition); ok && f.Name != nil {
			c.fragments[f.Name.Value] = f
		}
	}

	for _, def := range doc.Definitions {
		if op, ok := def.(*ast.OperationDefinition); ok {
			c.walk(op.SelectionSet, 1, 0, true)
		}
	}

	return c.rootFields <= maxRootFields && !c.aborted && c.cost <= maxQueryCost
}

func validateQueryDepth(query string, maxDepth int) bool {
	depth := 0
	maxReached := 0
	for _, char := range query {
		if char == '{' {
			depth++
			if depth > maxReached {
				maxReached = depth
			}
		} else if char == '}' {
			depth--
		}
	}
	return maxReached <= maxDepth
}

func (c *queryComplexity) walk(set *ast.SelectionSet, mult, depth int, root bool) {
	if set == nil || c.aborted {
		return
	}
	for _, sel := range set.Selections {
		c.steps++
		if c.steps > maxCostWalkSteps || depth > maxCostWalkDepth || c.cost > maxQueryCost {
			c.aborted = true
			return
		}

		switch s := sel.(type) {
		case *ast.Field:
			if root {
				c.rootFields++
			}
			if s.SelectionSet == nil {
				continue
			}
			c.cost += mult

			child := mult * listSize(s)
			if child > maxQueryCost {
				child = maxQueryCost + 1
			}
			c.walk(s.SelectionSet, child, depth+1, false)

		case *ast.InlineFragment:
			c.walk(s.SelectionSet, mult, depth+1, root)

		case *ast.FragmentSpread:
			if s.Name == nil {
				continue
			}
			if f, ok := c.fragments[s.Name.Value]; ok {
				c.walk(f.SelectionSet, mult, depth+1, root)
			}
		}
	}
}

func listSize(f *ast.Field) int {
	for _, arg := range f.Arguments {
		if arg.Name == nil || arg.Name.Value != "limit" {
			continue
		}
		v, ok := arg.Value.(*ast.IntValue)
		if !ok {
			return maxLimit
		}
		n, err := strconv.Atoi(v.Value)
		if err != nil || n > maxLimit {
			return maxLimit
		}
		if n < 1 {
			return 1
		}
		return n
	}
	return 1
}
