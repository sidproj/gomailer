package controller

import (
	"html/template"
	"net/http"
)


type HomeData struct{
	Title string
	Loop []string
}


func HomeController(w http.ResponseWriter, r *http.Request) {
	loop := []string{"arceus","zekrom","mewtwo"}
	p := HomeData{Title: "Pokemon",Loop: loop}
	t,_ := template.ParseFiles("views\\index.html")
	t.Execute(w,p)
}