package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/martinyonathann/restapi_golang_postgres/api/auth"
	"github.com/martinyonathann/restapi_golang_postgres/api/models"
	"github.com/martinyonathann/restapi_golang_postgres/api/responses"
	"github.com/martinyonathann/restapi_golang_postgres/api/utils/formaterror"
)

//CreateUser for create User
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		response := responses.Message("01", false, formattedError.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())
		response := responses.Message("01", false, formattedError.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)
}

//GetUsers for get User
func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

//UpdateUser for update user
func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	if tokenID != uint32(uid) {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		response := responses.Message("01", false, formattedError.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	responses.JSON(w, http.StatusOK, updatedUser)
}

//DeleteUser is fuction for delete user
func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		response := responses.Message("01", false, errors.New("Unauthorized").Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		response := responses.Message("01", false, errors.New(http.StatusText(http.StatusUnauthorized)).Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")

}
