package httpserver

import (
	"contacts/state"
	utils "contacts/utils"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func handleRefreshToken(s *state.State) http.HandlerFunc {

	type requestPayload struct {
		RefreshToken string `json:"refresh_token"`
	}

	type responsePayload struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	return func(w http.ResponseWriter, req *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				log.Error().Msgf("Recovered from panic: %v", r)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		refreshRequest := requestPayload{}

		err := json.NewDecoder(req.Body).Decode(&refreshRequest)
		if err != nil || refreshRequest.RefreshToken == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(refreshRequest.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Cfg.SecretKey), nil
		})

		if err != nil || !token.Valid {
			log.Error().Err(err).Msg("Invalid refresh token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.FromString(claims.Subject)

		if err != nil {
			log.Error().Err(err).Msg("Error parsing UUID")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		accessToken, err := utils.GenerateJWT(userID, utils.ScopeAuthentication, s.Cfg.SecretKey, 2*time.Hour)
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate access token")
			http.Error(w, "Error generating access token", http.StatusInternalServerError)
			return
		}

		newRefreshToken, err := utils.GenerateRefreshToken(claims.Subject, s.Cfg.SecretKey)
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate refresh token")
			http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
			return
		}

		tokenResponse := responsePayload{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
		}

		_ = loginSuccess.WriteToResponse(w, tokenResponse)
		return
	}
}
