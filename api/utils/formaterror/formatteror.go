//Create Custom Error Handling
package formaterror

import (
	"errors"
	"log"
	"strings"
)

func Message(rc string, status bool, message error) map[string]interface{} {
	var valueStatus string
	if status == false {
		valueStatus = "Failed"
	} else {
		valueStatus = "Success"
	}
	return map[string]interface{}{"rc": rc, "detail": valueStatus, "message": message}
}

func FormatError(err string) error {

	log.Printf("isi response error = %s", err) //debuging
	
	if strings.Contains(err, "nickname") {
		return errors.New("Nickname Already Taken")
	}
	if strings.Contains(err, "email") {
		return errors.New("Email Already Taken")
	}
	if strings.Contains(err, "title") {
		return errors.New("Title Already Taken")
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Password")
	}
	return errors.New("Incorrect Details")

}
