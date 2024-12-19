package controller

import (
	"html/template"
	"net/http"
)

func CreateTemplateControllerGET(w http.ResponseWriter,r * http.Request){
	t,_ := template.ParseFiles("views\\createTemplate.html")
	t.Execute(w,"")
}