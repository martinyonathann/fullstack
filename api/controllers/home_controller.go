package controllers

import (
	"net/http"

	"github.com/martinyonathann/fullstack/api/responses"
)

func (server *server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to this awesome API")
}
