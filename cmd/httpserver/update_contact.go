package httpserver

import (
	"contacts/repository"
	"contacts/state"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"net/http"
)

func PatchContactHandler(s *state.State) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		contactID := chi.URLParam(r, "id")
		uuidContactID, err := uuid.FromString(contactID)
		if err != nil {
			http.Error(w, "Invalid contact ID", http.StatusBadRequest)
			return
		}

		contact, err := s.Repository.GetContactByID(r.Context(), uuidContactID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		requestPayload := ContactRequestPayload{}
		if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if requestPayload.Phone != "" {
			contact.Phone = requestPayload.Phone
		}
		if requestPayload.Street != "" {
			contact.Street = requestPayload.Street
		}
		if requestPayload.City != "" {
			contact.City = requestPayload.City
		}
		if requestPayload.State != "" {
			contact.State = requestPayload.State
		}
		if requestPayload.ZipCode != "" {
			contact.ZipCode = requestPayload.ZipCode
		}
		if requestPayload.Country != "" {
			contact.Country = requestPayload.Country
		}

		updatedContact := repository.Contact{
			Phone:   contact.Phone,
			Street:  contact.Street,
			City:    contact.City,
			State:   contact.State,
			ZipCode: contact.ZipCode,
			Country: contact.Country,
		}

		err = s.Repository.PatchContact(r.Context(), uuidContactID, &updatedContact)
		if err != nil {
			http.Error(w, "Failed to update contact", http.StatusInternalServerError)
			return
		}
		response := ContactResponse{
			ID:      contactID,
			Phone:   contact.Phone,
			Street:  contact.Street,
			City:    contact.City,
			State:   contact.State,
			ZipCode: contact.ZipCode,
			Country: contact.Country,
		}

		_ = ContactUpdated.WriteToResponse(w, response)
		return
	}
}
