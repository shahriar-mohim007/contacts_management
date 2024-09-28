package httpserver

import (
	"contacts/state"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
)

type ContactResponse struct {
	ID      string `json:"id"`
	Phone   string `json:"phone"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

type ContactsResponse struct {
	Contacts []ContactResponse `json:"contacts"`
}

func GetAllContactsHandler(s *state.State) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			r := recover()
			if r != nil {
				log.Error().Msgf("Recovered from panic: %v", r)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		userID, _ := GetUserIDFromContext(r.Context())
		uuID, err := uuid.FromString(userID)

		if err != nil {
			log.Error().Err(err).Msg("Error parsing UUID")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		contacts, err := s.Repository.GetAllContacts(r.Context(), uuID)
		if err != nil {
			log.Error().Err(err).Msg("Error fetching contacts")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var contactResponses []ContactResponse
		for _, contact := range contacts {
			contactResponses = append(contactResponses, ContactResponse{
				ID:      contact.ID.String(),
				Phone:   contact.Phone,
				Street:  contact.Street,
				City:    contact.City,
				State:   contact.State,
				ZipCode: contact.ZipCode,
				Country: contact.Country,
			})
		}

		response := ContactsResponse{
			Contacts: contactResponses,
		}

		_ = ContactRetrieved.WriteToResponse(w, response)
		return
	}
}
