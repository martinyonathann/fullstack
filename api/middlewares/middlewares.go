package middlewares

import (
	"net/http"

	"github.com/martinyonathann/restapi_golang_postgres/api/auth"
	"github.com/martinyonathann/restapi_golang_postgres/api/responses"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			response := responses.Message("01", false, err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			responses.Respond(w, response)
			return
		}
		next(w, r)
	}
}
