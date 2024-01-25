package interceptors

import (
	"api/constants"
	"api/models"
	"context"
	"errors"
)

const EMPTY = ""

var ErrUnableToParseAccount = errors.New("unable to parse authenticated account")
var ErrUnableToParseToken = errors.New("unable to parse authenticated account token")

// GetAuthenticatedAccount returns the authenticated account from the context.
func GetAuthenticatedAccount(ctx context.Context) (*models.User, error) {
	account, ok := ctx.Value(constants.AuthenticatedAccountContextKey).(*models.User)
	if !ok {
		return nil, ErrUnableToParseAccount
	}

	return account, nil
}

func GetAuthenticatedAccountToken(ctx context.Context) (token string, err error) {
	var ok bool

	token, ok = ctx.Value(constants.AuthenticatedSessionTokenContextKey).(string)
	if !ok {
		return EMPTY, ErrUnableToParseToken
	}

	return token, nil
}
