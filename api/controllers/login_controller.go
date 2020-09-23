package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/martinyonathann/restapi_golang_postgres/api/auth"
	"github.com/martinyonathann/restapi_golang_postgres/api/models"
	"github.com/martinyonathann/restapi_golang_postgres/api/responses"
	"github.com/martinyonathann/restapi_golang_postgres/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	// buf, bodyErr := ioutil.ReadAll(r.Body)
	// if bodyErr != nil {
	// 	log.Print("bodyErr ", bodyErr.Error())
	// 	http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	// log.Printf("BODY: %q", rdr1)

	response := make(map[string]interface{})
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		// responses.ERROR(w, http.StatusUnprocessableEntity, err)
		response = responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		// responses.ERROR(w, http.StatusUnprocessableEntity, err)
		response = responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		// responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		response = responses.Message("01", false, formattedError.Error())
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	resp := map[string]interface{}{"rc": http.StatusOK, "detail": "Success", "message": "Success", "token": "Bearer " + token}
	w.WriteHeader(http.StatusOK)
	responses.Respond(w, resp)
}
func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
