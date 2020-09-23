package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

//JSON for handle JSON
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		// fmt.Fprint(w, "%s", err.Error())
	}
}

//Message for message error
func Message(rc string, status bool, message string) map[string]interface{} {
	var valueStatus string
	if status == false {
		valueStatus = "Failed"
	} else {
		valueStatus = "Success"
	}
	resp := map[string]interface{}{"rc": rc, "detail": valueStatus, "message": message}
	log.Printf("isi response : rc %s message %s ", rc, message)
	return resp
}

//Respond for response json
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
