package controllers

import (
	"net/http"

	"github.com/martinyonathann/restapi_golang_postgres/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	response := responses.Message("00", true, "Welcome to this awesome API")
	w.WriteHeader(http.StatusForbidden)
	w.Header().Add("Content-Type", "application/json")
	responses.Respond(w, response)
	return
}
