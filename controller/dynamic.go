package controller

import (
	"encoding/json"
	"fmt"
	"gomailer/router"
	"net/http"
)

func DynamicControllerGET(w http.ResponseWriter, r *http.Request) {
	params := router.GetPathParams(r)
	id := params["id"]
	fmt.Fprintf(w, "ID: %s", id)

	json.NewEncoder(w).Encode(params)
}

func Dynamic1ControllerGET(w http.ResponseWriter, r *http.Request) {
	params := router.GetPathParams(r)
	fmt.Printf("testing : ID: %+v",params)

	json.NewEncoder(w).Encode(params)
}