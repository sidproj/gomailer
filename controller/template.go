package controller

import (
	"encoding/json"
	"gomailer/models"
	"gomailer/utils"
	"html/template"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	allTemplateView    = "allTemplates.html"
	createTemplateView = "createTemplate.html"
)

type CreateTemplateRequest struct {
	TemplateContent   string
	TemplateName      string
	TemplateVariables []string
	ID                string
}

type DeleteTemplateRequest struct {
	TemplateID string
}

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description,omitempty"`
}

func handleErrorRedirect(w http.ResponseWriter, r *http.Request, err error, path string) {
	log.Println(err)
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func TemplateControllerGET(w http.ResponseWriter, r *http.Request) {
	templateModal, err := models.GetTemplateModel()

	if err != nil {
		handleErrorRedirect(w, r, err, "/template")
		return
	}

	filter := bson.M{}

	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}})

	templates, err := templateModal.Find(filter, opts)

	if err != nil {
		handleErrorRedirect(w, r, err, "/template")
		return
	}

	templateList := []map[string]interface{}{}

	for _, val := range templates {
		templateList = append(templateList, map[string]interface{}{
			"id":      val.ID.Hex(),
			"name":    val.Name,
			"content": val.TemplateContent,
		})
	}

	jsonBytes, err := json.Marshal(templateList)
	if err != nil {
		http.Error(w, "Failed to serialize templates", http.StatusInternalServerError)
		return
	}

	templateData := map[string]interface{}{
		"templateListJSON": template.JS(jsonBytes),
	}
	utils.RenderTemplate(w, allTemplateView, templateData)
}

func CreateTemplateControllerGET(w http.ResponseWriter, r *http.Request) {
	// csrf token generation
	token, err := utils.GenerateCSRFToken()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	utils.SetCSRFCookie(w, token)

	templateData := map[string]interface{}{
		"CsrfToken": token,
	}
	utils.RenderTemplate(w, createTemplateView, templateData)
}

func CreateTemplateControllerPOST(w http.ResponseWriter, r *http.Request) {

	templateData := CreateTemplateRequest{
		TemplateContent:   utils.TrimSpaces(r.FormValue("templateContent")),
		TemplateVariables: []string{},
		TemplateName:      utils.TrimSpaces(r.FormValue("templateName")),
	}
	w.Header().Set("Content-Type", "application/json")

	// csrf token verification
	if !utils.VerifyCSRFToken(r) {
		token, err := utils.GenerateCSRFToken()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		utils.SetCSRFCookie(w, token)
		templateData := map[string]interface{}{
			"CsrfToken": token,
			"error":     "CSRF token mismatch. Please reload the page",
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(templateData)
		return
	}

	variables := strings.Split(r.FormValue("templateVariables"), ",")

	for _, vars := range variables {
		if len(vars) > 0 {
			templateData.TemplateVariables = append(templateData.TemplateVariables, vars)
		}
	}

	userID, err := primitive.ObjectIDFromHex(r.Header.Get("user_id"))

	if err != nil {
		errorMap := map[string]string{
			"error":    "no user found!",
			"redirect": "/login",
		}

		json.NewEncoder(w).Encode(errorMap)
		return
	}

	templateModel, err := models.GetTemplateModel()

	if err != nil {
		errorMap := map[string]string{
			"error": "error while getting template modal",
		}

		json.NewEncoder(w).Encode(errorMap)
		return
	}

	template := models.TemplateSchema{
		UserID:            userID,
		Name:              templateData.TemplateName,
		TemplateContent:   utils.TrimSpaces(templateData.TemplateContent),
		TemplateVariables: templateData.TemplateVariables,
	}

	err = template.Validate()

	if err != nil {
		errorMap := map[string]string{
			"error":       "error while saving template",
			"description": err.Error(),
		}
		log.Println(err.Error())
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	if err := templateModel.Save(&template); err != nil {
		errMsg := err.Error()
		if mongo.IsDuplicateKeyError(err) {
			errMsg = "Template with the title already exists"
		}
		errorMap := map[string]string{
			"error": errMsg,
		}
		log.Println(err.Error())
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	w.WriteHeader(http.StatusCreated)
	tempMap := map[string]string{
		"message":  "success",
		"redirect": "/template/edit?template_id=" + template.ID.Hex(),
	}
	json.NewEncoder(w).Encode(tempMap)
}
