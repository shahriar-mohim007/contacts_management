package httpserver

import (
	"contacts/state"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"net/http"
)

func DeleteContactHandler(s *state.State) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")

		contactID, err := uuid.FromString(id)
		if err != nil {
			http.Error(w, "Invalid contact ID", http.StatusBadRequest)
			return
		}

		err = s.Repository.DeleteContactByID(r.Context(), contactID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Contact not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}
}
