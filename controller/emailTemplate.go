package controller

import (
	"encoding/json"
	"gomailer/models"
	"gomailer/utils"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	allTemplateView = filepath.Join("views","templates.html")
	createTemplateView = filepath.Join("views","createTemplate.html")
)

type CreateTemplateRequest struct{
	TemplateContent string
	TemplateName string
	TemplateVariables []string
	ID string
}

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description,omitempty"`
}

func handleErrorRedirect(w http.ResponseWriter, r *http.Request, err error, path string) {
	log.Println(err)
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func TemplateControllerGET(w http.ResponseWriter,r * http.Request){
	templateModal,err := models.GetTemplateModel()
	
	if(err!=nil){
		handleErrorRedirect(w,r,err,"/template")
		return
	}

	filter := bson.M{}

	templates,err := templateModal.Find(filter)

	if(err!=nil){
		handleErrorRedirect(w,r,err,"/template")
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
	utils.RenderTemplate(w,allTemplateView,templateData)
}

func CreateTemplateControllerGET(w http.ResponseWriter,r * http.Request){
	// csrf token generation
	token, err := utils.GenerateCSRFToken()
	if(err!=nil){
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	utils.SetCSRFCookie(w, token)
	
	templateData := map[string]interface{}{
		"CsrfToken": token,
	}
	utils.RenderTemplate(w,createTemplateView,templateData)
}

func CreateTemplateControllerPOST(w http.ResponseWriter,r* http.Request){

	templateData := CreateTemplateRequest{
		TemplateContent: utils.TrimSpaces(r.FormValue("templateContent")),
		TemplateVariables: []string{},
		TemplateName: utils.TrimSpaces(r.FormValue("templateName")),
	}
	w.Header().Set("Content-Type", "application/json")

	// csrf token verification
	if !utils.VerifyCSRFToken(r) {	
		token, err := utils.GenerateCSRFToken()
		if(err!=nil){
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		utils.SetCSRFCookie(w, token)
		templateData := map[string]interface{}{
			"CsrfToken": token,
			"error":"CSRF token mismatch. Please reload the page",
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(templateData)
		return
	}

	
	variables := strings.Split(r.FormValue("templateVariables"),",")

	for _,vars:=range variables {
		if(len(vars)>0){
			templateData.TemplateVariables = append(templateData.TemplateVariables,vars)
		}
	}

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
		TemplateContent: utils.TrimSpaces(templateData.TemplateContent),
		TemplateVariables: templateData.TemplateVariables,
	}

	err = template.Validate()

	if(err!=nil){
		errorMap := map[string]string{
			"error":"error while saving template",
			"description":err.Error(),
		}
        log.Println(err.Error())
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	if err:=templateModel.Save(&template);err!=nil{
		errMsg := err.Error()
		if mongo.IsDuplicateKeyError(err){
			errMsg = "Template with the title already exists"
		}
		errorMap := map[string]string{
			"error":errMsg,
		}
        log.Println(err.Error())
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	w.WriteHeader(http.StatusCreated)
	tempMap := map[string]string{
		"message":"success",
	}
    json.NewEncoder(w).Encode(tempMap)
}

func EditTemplateControllerGET(w http.ResponseWriter,r* http.Request){

	templateID := r.URL.Query().Get("template_id")

	if(templateID == ""){
		http.Redirect(w,r,"/template",http.StatusSeeOther)
		return
	}

	// csrf token generation
	token, err := utils.GenerateCSRFToken()
	if(err!=nil){
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	utils.SetCSRFCookie(w, token)

	templateModel,err := models.GetTemplateModel()

	if(err!=nil){
		handleErrorRedirect(w,r,err,"/template")
		return
	}

	user_id,err := primitive.ObjectIDFromHex(r.Header.Get("user_id"))
	if (err != nil){
		handleErrorRedirect(w,r,err,"/template")
		return
	}

	template_id,err := primitive.ObjectIDFromHex(templateID)
	if (err != nil){
		handleErrorRedirect(w,r,err,"/template")
		return
	}

	filter := bson.M{"_id":template_id,"user_id":user_id,}

	data,err := templateModel.Find(filter)

	if(err!=nil){
		handleErrorRedirect(w,r,err,"/template")
		return
	}

	templateData := map[string]interface{}{
		"templateContent":data[0].TemplateContent,
		"templateName":data[0].Name,
		"templateVariables":strings.Join(data[0].TemplateVariables,","),
		"CsrfToken":token,
	}
	utils.RenderTemplate(w,createTemplateView,templateData)
}

func EditTemplateControllerPOST(w http.ResponseWriter,r* http.Request){
	
	templateID := r.URL.Query().Get("template_id")

	if(templateID == ""){
		http.Redirect(w,r,"/template",http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if !utils.VerifyCSRFToken(r) {
		errorMap := map[string]interface{}{
			"error":"CSRF token mismatch. Please reload the page",
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	templateData := CreateTemplateRequest{
		TemplateContent: utils.TrimSpaces(r.FormValue("templateContent")),
		TemplateName: utils.TrimSpaces(r.FormValue("templateName")),
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
        log.Println(err.Error())
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
        log.Println(err.Error())
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	tempMap := make(map[string]interface{})

	tempMap["content"] = oldTemplate.TemplateContent
	tempMap["variables"] = oldTemplate.TemplateVariables

    json.NewEncoder(w).Encode(tempMap)

}