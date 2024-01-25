package middlewares

import (
	"api/constants"
	"api/models"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

var (
	ErrUnauthorized = errors.New("unauthorized: Sorry, you must be authenticated/logged in to continue")
	ErrTokenMissing = errors.New("token missing in Authorization header")
	ErrInvalidToken = errors.New("invalid token")
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func (mw *SystemMiddleware) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	response, err := mw.userService.VerifyAuthToken(ctx, token)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return response, nil
}

// AuthMiddleware is a middleware that checks if the user is authenticated.
func (mw *SystemMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw.logger.Infof("AuthMiddleware: Request URL - %s", r.URL.Path)
		if _, isWhitelisted := constants.WhitelistedRoutes[r.URL.Path]; isWhitelisted {
			next.ServeHTTP(w, r)
			return
		}
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			respondWithError(w, ErrUnauthorized)
			return
		}
		sp := strings.Split(authorization, " ")
		if len(sp) != 2 {
			respondWithError(w, ErrInvalidToken)
			return
		}
		token := sp[1]
		ctx, err := mw.ValidateHeaders(r.Context(), token, len(token) > 0)

		if err != nil {
			mw.logger.WithError(err).Error("Authentication failed")
			respondWithError(w, ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// respondWithError responds with an error message.
func respondWithError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	result := ErrorResponse{
		Message: err.Error(),
	}
	if encodeErr := json.NewEncoder(w).Encode(result); encodeErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
