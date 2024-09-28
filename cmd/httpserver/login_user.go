package httpserver

import (
	"contacts/state"
	utils "contacts/utils"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func handleLogin(s *state.State) http.HandlerFunc {

	type requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type responsePayload struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Msgf("Recovered from panic: %v", r)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		request := requestPayload{}
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil || request.Email == "" || request.Password == "" {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		user, err := s.Repository.GetUserByEmail(req.Context(), request.Email)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		if !utils.CheckPasswordHash(user.Password, request.Password) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		if !user.IsActive {
			http.Error(w, "User not active", http.StatusUnauthorized)
			return
		}

		ttl := 2 * time.Hour

		accessToken, err := utils.GenerateJWT(user.ID, utils.ScopeAuthentication, s.Cfg.SecretKey, ttl)
		if err != nil {
			log.Error().Err(err).Msg("Error generating access token")
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := utils.GenerateRefreshToken(user.ID.String(), s.Cfg.SecretKey)
		if err != nil {
			http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
			return
		}

		response := responsePayload{
			Token:        accessToken,
			RefreshToken: refreshToken,
		}
		_ = loginSuccess.WriteToResponse(w, response)

		return

	}
}
