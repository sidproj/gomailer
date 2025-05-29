package controller

import (
	"encoding/json"
	"fmt"
	"gomailer/models"
	"gomailer/utils"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TemplateData struct{
	TemplateContent string
	TemplateName string
	TemplateVariables []string
	ID string
}

func TemplateControllerGET(w http.ResponseWriter,r * http.Request){
	templateModal,err := models.GetTemplateModel()
	
	if(err!=nil){
		fmt.Println(err.Error())
		http.Redirect(w,r,"/template",http.StatusSeeOther)
		return
	}

	filter := bson.M{}

	templates,err := templateModal.Find(filter)

	if(err!=nil){
		fmt.Println(err.Error())
		http.Redirect(w,r,"/template",http.StatusSeeOther)
		return
	}

	templateData := map[string]interface{}{
		"templateList":[]interface{}{},
	}

	for _,val := range templates{
		newTemplate := map[string]interface{}{
			"id":val.ID.Hex(),
			"name":val.Name,
		}

		templateData["templateList"] = append(
			templateData["templateList"].([]interface{}),
			newTemplate)
	}
	utils.RenderTemplate(w,"views\\templates.html",templateData)
}

func CreateTemplateControllerGET(w http.ResponseWriter,r * http.Request){
	templateData := map[string]interface{}{
	}
	utils.RenderTemplate(w,"views\\createTemplate.html",templateData)
}

func CreateTemplateControllerPOST(w http.ResponseWriter,r* http.Request){

	templateData := TemplateData{
		TemplateContent: strings.Trim(r.FormValue("templateContent")," "),
		TemplateVariables: []string{},
		TemplateName: strings.Trim(r.FormValue("templateName")," "),
	}
	
	variables := strings.Split(r.FormValue("templateVariables"),",")

	for _,vars:=range variables {
		if(len(vars)>0){
			templateData.TemplateVariables = append(templateData.TemplateVariables,vars)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	userID,err := primitive.ObjectIDFromHex(r.Header.Get("user_id"))

	if(err!=nil){
		errorMap := map[string]string{
			"error":"no user found!",
			"redirect":"/login",
		}

		json.NewEncoder(w).Encode(errorMap)
		return
	}

	templateModel,err := models.GetTemplateModel()

	if(err!=nil){
		errorMap := map[string]string{
			"error":"error while getting template modal",
		}

		json.NewEncoder(w).Encode(errorMap)
		return
	}

	template := models.TemplateSchema{
		UserID : userID,
		Name: templateData.TemplateName,
		TemplateContent: strings.Trim(templateData.TemplateContent," "),
		TemplateVariables: templateData.TemplateVariables,
	}

	err = template.Validate()

	if(err!=nil){
		errorMap := map[string]string{
			"error":"error while saving template",
			"description":err.Error(),
		}
        fmt.Print(err.Error())
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	if err:=templateModel.Save(&template);err!=nil{
		errorMap := map[string]string{
			"error":"error while saving template",
			"description":err.Error(),
		}
        fmt.Print(err.Error())
		json.NewEncoder(w).Encode(errorMap)
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
		return
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
		"templateName":data[0].Name,
		"templateVariables":strings.Join(data[0].TemplateVariables,","),
	}
	utils.RenderTemplate(w,"views\\createTemplate.html",templateData)
}

func EditTemplateControllerPOST(w http.ResponseWriter,r* http.Request){
	
	templateID := r.URL.Query().Get("template_id")

	if(templateID == ""){
		http.Redirect(w,r,"/template",http.StatusSeeOther)
		return
	}

	templateData := TemplateData{
		TemplateContent: strings.Trim(r.FormValue("templateContent")," "),
		TemplateName: strings.Trim(r.FormValue("templateName")," "),
		TemplateVariables: []string{},
	}

	variables := strings.Split(r.FormValue("templateVariables"),",")

	for _,vars:=range variables{
		if(len(vars)>0){
			templateData.TemplateVariables = append(templateData.TemplateVariables, vars)
		}
	}

	templateModel,err := models.GetTemplateModel()

	if(err!=nil){
		errorMap := map[string]string{
			"error":"error while getting template modal",
		}

		json.NewEncoder(w).Encode(errorMap)
		return
	}

	oldTemplate,err := templateModel.FindById(templateID)

	if(err != nil){
		errorMap := map[string]string{
			"error":"error while finding template",
			"description":err.Error(),
		}
        fmt.Print(err.Error())
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	oldTemplate.Name = templateData.TemplateName
	oldTemplate.TemplateContent = templateData.TemplateContent
	oldTemplate.TemplateVariables = templateData.TemplateVariables

	if err:= templateModel.Save(&oldTemplate);err!=nil{
		errorMap := map[string]string{
			"error":"error while updating template",
			"description":err.Error(),
		}
        fmt.Print(err.Error())
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	tempMap := make(map[string]interface{})

	tempMap["content"] = oldTemplate.TemplateContent
	tempMap["variables"] = oldTemplate.TemplateVariables

    json.NewEncoder(w).Encode(tempMap)

}