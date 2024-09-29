package httpserver

import (
	"contacts/state"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func GetContactByIDHandler(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")

		contactID, err := uuid.FromString(id)
		if err != nil {
			http.Error(w, "Invalid contact ID", http.StatusBadRequest)
			return
		}

		contact, err := s.Repository.GetContactByID(r.Context(), contactID)
		if err != nil {
			if strings.Contains(err.Error(), "no contact found") {
				http.Error(w, "Contact not found", http.StatusNotFound)
			} else {
				log.Error().Err(err).Msg("Error fetching contact")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		_ = ContactRetrieved.WriteToResponse(w, contact)
		return
	}
}
