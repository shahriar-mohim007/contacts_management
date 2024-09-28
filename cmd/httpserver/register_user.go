package httpserver

import (
	"contacts/repository"
	"contacts/state"
	utils "contacts/utils"
	"database/sql"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func handleRegisterUser(s *state.State) http.HandlerFunc {
	type requestPayload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type responsePayload struct {
		ActivateToken string `json:"activate_token"`
	}

	return func(w http.ResponseWriter, req *http.Request) {
		request := requestPayload{}
		ctx := req.Context()
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil || request.Name == "" || request.Email == "" || request.Password == "" {
			log.Fatal().Err(err).Msgf("Invalid input: %v", err)
			_ = ValidDataNotFound.WriteToResponse(w, nil)
			return
		}
		userDto, err := s.Repository.GetUserByEmail(ctx, request.Email)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Debug().Msgf("User with email %s not found", request.Email)
			} else {
				log.Error().Err(err).Msg("Error fetching user by email")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		if userDto != nil {
			log.Debug().Msgf("User already exists: %v", request.Email)
			_ = UserAlreadyExist.WriteToResponse(w, nil)
			return
		}

		passwordHash, err := utils.HashPassword(request.Password)

		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to hash password")
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.NewV4()
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		user := repository.User{
			ID:       userID,
			Name:     request.Name,
			Email:    request.Email,
			Password: passwordHash,
			IsActive: false,
		}

		if err := s.Repository.CreateUser(ctx, &user); err != nil {
			log.Fatal().Err(err).Msgf("Failed to create user")
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return

		}

		ttl := 2 * time.Hour

		token, err := utils.GenerateJWT(user.ID, utils.ScopeActivation, s.Cfg.SecretKey, ttl)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to activation token")
			http.Error(w, "Error creating activation token", http.StatusInternalServerError)
			return
		}

		response := responsePayload{
			ActivateToken: token,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)

		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

	}
}
