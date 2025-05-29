package controller

import (
	"gomailer/utils"
	"net/http"
)


type HomeData struct{
	Title string
	Loop []string
}


func HomeController(w http.ResponseWriter, r *http.Request) {
	loop := []string{"arceus","zekrom","mewtwo"}
	p := HomeData{Title: "Pokemon",Loop: loop}
	utils.RenderTemplate(w,"views\\index.html",p)
}