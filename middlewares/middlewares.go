package middlewares

import (
	"api/constants"
	"api/services"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrTokenRequired  = errors.New("token is required")
	ErrSessionInvalid = errors.New("session is invalid")
)

const EMPTY = ""

type SystemMiddleware struct {
	userService *services.UserService
	logger      *logrus.Entry
}

func NewSystemMiddleware(
	userService *services.UserService,
	logger *logrus.Logger,
) *SystemMiddleware {
	return &SystemMiddleware{
		userService: userService,
		logger:      logger.WithField("component", "SystemMiddleware"),
	}
}

// ValidateHeaders validates the headers of the request.
func (mw *SystemMiddleware) ValidateHeaders(ctx context.Context, token string, ensureTokenValidation bool) (context.Context, error) {
	if ensureTokenValidation {
		if token == EMPTY {
			return ctx, ErrTokenRequired
		}
	}

	profile, err := mw.ValidateToken(ctx, token)
	if err != nil {
		mw.logger.WithError(err).WithField("header.token", token).Error("failed to fetch account from token")
		return ctx, ErrSessionInvalid
	}

	ctx = context.WithValue(ctx, constants.AuthenticatedAccountContextKey, profile)
	ctx = context.WithValue(ctx, constants.AuthenticatedSessionTokenContextKey, token)
	return ctx, nil
}
