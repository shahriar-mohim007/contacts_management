package httpserver

import (
	"contacts/repository"
	"contacts/state"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
)

type ContactRequestPayload struct {
	Phone   string `json:"phone" validate:"required"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

func CreateContactHandler(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestPayload := ContactRequestPayload{}
		if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		userID, _ := GetUserIDFromContext(r.Context())
		uuID, err := uuid.FromString(userID)

		if err != nil {
			log.Error().Err(err).Msg("Error parsing UUID")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		ID, err := uuid.NewV4()
		if err != nil {
			http.Error(w, "Error Generating Id", http.StatusInternalServerError)
			return
		}

		contact := repository.Contact{
			ID:      ID,
			UserID:  uuID,
			Phone:   requestPayload.Phone,
			Street:  requestPayload.Street,
			City:    requestPayload.City,
			State:   requestPayload.State,
			ZipCode: requestPayload.ZipCode,
			Country: requestPayload.Country,
		}

		if err := s.Repository.CreateContact(r.Context(), &contact); err != nil {
			log.Printf("Error creating contact: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		_ = ContactCreated.WriteToResponse(w, contact)

		return
	}
}
