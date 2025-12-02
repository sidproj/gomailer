package controller

import (
	"gomailer/utils"
	"net/http"
)

var (
	aboutHtml = "about.html"
)

func AboutController(w http.ResponseWriter, _ *http.Request) {
	utils.RenderTemplate(w, aboutHtml, nil)
}
