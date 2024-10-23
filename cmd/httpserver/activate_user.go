package httpserver

import (
	"contacts/state"
	utils "contacts/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

func HandleActivateUser(s *state.State) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			r := recover()
			if r != nil {
				log.Error().Msgf("Recovered from panic: %v", r)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		tokenString := r.URL.Query().Get("token")
		ctx := r.Context()

		if tokenString == "" {
			http.Error(w, "Token is required", http.StatusBadRequest)
			return
		}

		var claims utils.Claims
		secretKey := s.Cfg.SecretKey

		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID := claims.UserID

		err = s.Repository.ActivateUserByID(ctx, userID)
		if err != nil {
			log.Error().Err(err).Msg("Error activating user")
			http.Error(w, "Error activating user", http.StatusInternalServerError)
			return
		}
		_ = UserActivated.WriteToResponse(w, nil)
		return

	}
}
