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

func (server *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	post.Prepare()
	err = post.Validate()
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		response := responses.Message("01", false, errors.New("Unauthorized").Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	if uid != post.AuthorID {
		response := responses.Message("01", false, errors.New("Unauthorized").Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	if uid != post.AuthorID {
		response := responses.Message("01", false, errors.New(http.StatusText(http.StatusUnauthorized)).Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		response := responses.Message("01", false, formattedError.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JSON(w, http.StatusCreated, postCreated)
}

func (server *Server) GetPosts(w http.ResponseWriter, r *http.Request) {

	post := models.Post{}

	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}

func (server *Server) GetPost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	post := models.Post{}

	postReceived, err := post.FindPostByID(server.DB, pid)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	responses.JSON(w, http.StatusOK, postReceived)
}
func (server *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//Check if the post id iss valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	//check if the auth token is valid and get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		response := responses.Message("01", false, errors.New("Unauthorized").Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}

	//check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		response := responses.Message("01", false, errors.New("Post not found").Error())
		w.WriteHeader(http.StatusNotFound)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	//if a user attempt to update a post not belonging to him
	if uid != post.AuthorID {
		response := responses.Message("01", false, errors.New("Unauthorized").Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	//read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	//start processing the request data
	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	//Also check if the request user id is equal to the one gotten from token
	if uid != postUpdate.AuthorID {
		response := responses.Message("01", false, errors.New("Unauthorized").Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	postUpdate.Prepare()
	err = postUpdate.Validate()
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	postUpdate.ID = post.ID //this is important to tell the model the post id to update, the other update field are set above
	postUpdated, err := postUpdate.UpdateAPost(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		response := responses.Message("01", false, formattedError.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	responses.JSON(w, http.StatusOK, postUpdated)
}

func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//is a valid post id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	// is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		response := responses.Message("01", false, errors.New("Unauthorized").Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}

	//check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		response := responses.Message("01", false, errors.New("Unauthorized").Error())
		w.WriteHeader(http.StatusNotFound)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	_, err = post.DeleteAPost(server.DB, pid, uid)
	if err != nil {
		response := responses.Message("01", false, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		responses.Respond(w, response)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
