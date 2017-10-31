package main

import (
	"encoding/json"
	"fmt"
	validator "github.com/gidsi/go-spaceapi-validator"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", validateHandler).
		Methods("POST")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		statusCode, returnBody := mapValidation(string(body))

		w.WriteHeader(statusCode)
		w.Write(returnBody)
	}
}

func mapValidation(document string) (int, []byte) {
	result, err := validator.Validate(document)
	if err != nil {
		return http.StatusBadRequest, nil
	}

	if !result.Valid {
		response, err := json.Marshal(result.Errors)
		if err != nil {
			fmt.Println(err)
			return http.StatusInternalServerError, nil
		}

		return http.StatusUnprocessableEntity, response
	}
	return http.StatusOK, nil
}
