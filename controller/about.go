package controller

import (
	"fmt"
	"net/http"
)

func AboutController(w http.ResponseWriter, _ *http.Request){
	fmt.Fprintf(w,"<h1>About page.</h1>")
}