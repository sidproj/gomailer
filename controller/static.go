package controller

import (
	"gomailer/utils"
	"net/http"
)

var (
	aboutHtml     = "about.html"
	privacyPolicy = "privacyPolicy.html"
)

func AboutController(w http.ResponseWriter, _ *http.Request) {
	utils.RenderTemplate(w, aboutHtml, nil)
}

func PrivacyPollicyController(w http.ResponseWriter, _ *http.Request) {
	utils.RenderTemplate(w, privacyPolicy, nil)
}
