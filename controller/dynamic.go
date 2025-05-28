package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	router "github.com/sidproj/grouter"
)

func DynamicControllerGET(w http.ResponseWriter, r *http.Request) {
	params := router.GetPathParams(r)
	json.NewEncoder(w).Encode(params)
}

func Dynamic1ControllerGET(w http.ResponseWriter, r *http.Request) {
	params := router.GetPathParams(r)
	fmt.Printf("testing : ID: %+v",params)

	json.NewEncoder(w).Encode(params)
}