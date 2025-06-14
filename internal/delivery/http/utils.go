package http

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func bindAndValidate(w http.ResponseWriter, r *http.Request, dst interface{}, validate *validator.Validate) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid payload")
		return false
	}
	if err := validate.Struct(dst); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			errs := make([]string, 0, len(verrs))
			for _, verr := range verrs {
				errs = append(errs, verr.Field()+": "+verr.ActualTag())
			}
			respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{"validation_errors": errs})
		} else {
			respondWithError(w, http.StatusBadRequest, "bad request")
		}
		return false
	}
	return true
}
