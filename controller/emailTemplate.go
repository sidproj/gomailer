package controller

import (
	"encoding/json"
	"fmt"
	"gomailer/models"
	"html/template"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TemplateData struct{
	templateContent string
	templateVariables []string
}

func CreateTemplateControllerGET(w http.ResponseWriter,r * http.Request){
	templateData := map[string]interface{}{
	}
	t,_ := template.ParseFiles("views\\createTemplate.html")
	t.Execute(w,templateData)
}

func CreateTemplateControllerPOST(w http.ResponseWriter,r* http.Request){


	templateData := TemplateData{
		templateContent: strings.Trim(r.FormValue("templateContent")," "),
		templateVariables: []string{},
	}
	
	variables := strings.Split(r.FormValue("templateVariables"),",")

	for _,vars:=range variables {
		if(len(vars)>0){
			templateData.templateVariables = append(templateData.templateVariables,vars)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	userID,err := primitive.ObjectIDFromHex(r.Header.Get("user_id"))

	templateError := sendError(w,"views\\login.html")

	if(err!=nil){
		templateError("template creation error: error while getting user_id from request header. "+err.Error())
		return
	}

	templateModel,err := models.GetTemplateModel()

	if(err!=nil){
		templateError("template creation error: error while getting template model. "+err.Error())
		return
	}

	template := models.TemplateSchema{
		UserID : userID,
		TemplateContent: templateData.templateContent,
		TemplateVariables: templateData.templateVariables,
	}

	fmt.Println(templateData.templateVariables)
	fmt.Println("variables",templateData.templateVariables)

	if err:=templateModel.Save(template);err!=nil{
		templateError("template creation error: error while saving template. "+err.Error())
		return
	}

	tempMap := make(map[string]interface{})

	tempMap["content"] = template.TemplateContent
	tempMap["variables"] = template.TemplateVariables

    json.NewEncoder(w).Encode(tempMap)
}

func EditTemplateControllerGET(w http.ResponseWriter,r* http.Request){

	templateID := r.URL.Query().Get("template_id")

	if(templateID == ""){
		http.Redirect(w,r,"/template",http.StatusSeeOther)
	}

	templateModel,err := models.GetTemplateModel()

	if(err!=nil){
		fmt.Println(err.Error())
		http.Redirect(w,r,"/template",http.StatusSeeOther)
		return
	}

	user_id,err := primitive.ObjectIDFromHex(r.Header.Get("user_id"))
	if (err != nil){
		fmt.Println(err.Error())
		http.Redirect(w,r,"/template",http.StatusSeeOther)
		return
	}

	template_id,err := primitive.ObjectIDFromHex(templateID)
	if (err != nil){
		fmt.Println(err.Error())
		http.Redirect(w,r,"/template",http.StatusSeeOther)
		return
	}

	filter := bson.M{"_id":template_id,"user_id":user_id,}

	data,err := templateModel.Find(filter)

	if(err!=nil){
		fmt.Println(err.Error())
		http.Redirect(w,r,"/template",http.StatusSeeOther)
		return
	}

	templateData := map[string]interface{}{
		"templateContent":data[0].TemplateContent,
	}
	t,_ := template.ParseFiles("views\\createTemplate.html")
	t.Execute(w,templateData)
}