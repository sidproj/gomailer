package controller

import (
	"gomailer/utils"
	"net/http"
)

var (
	homeHTML = "home.html"
)

type HomeData struct {
	Title string
	Loop  []string
}

func HomeController(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, homeHTML, nil)
}
