package api

import (
	"log/slog"
	"strings"

	"github.com/graphql-go/graphql/gqlerrors"
)

func sanitizeErrors(logger *slog.Logger, errs []gqlerrors.FormattedError) {
	for i, err := range errs {
		logger.Error("GraphQL Error", "message", err.Message)

		switch {
		case strings.Contains(err.Message, "Cannot query field"):
			errs[i].Message = "Invalid field requested: " + err.Message
		case strings.Contains(err.Message, "Unknown argument"):
			errs[i].Message = "Invalid argument provided: " + err.Message
		case isClientError(err.Message):
		default:
			errs[i].Message = "An internal server error occurred while processing the request"
		}
	}
}

func isClientError(msg string) bool {
	for _, p := range []string{"cannot be resolved", "Variable", "was not provided", "Argument \""} {
		if strings.Contains(msg, p) {
			return true
		}
	}
	return false
}
