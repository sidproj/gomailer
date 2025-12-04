package controller

import (
	"encoding/json"
	"gomailer/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteTemplateControllerDELETE(w http.ResponseWriter, r *http.Request) {

	templateID := r.URL.Query().Get("template_id")

	if templateID == "" {
		sendJsonError(w, "Invalid request data", "No template id provided")
		return
	}

	user_id, err := primitive.ObjectIDFromHex(r.Header.Get("user_id"))

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

	template_id, err := primitive.ObjectIDFromHex(templateID)

	if err != nil {
		sendJsonError(w, "Invalid request data", "Invalid template id provided")
		return
	}

	filter := bson.M{
		"_id":     template_id,
		"user_id": user_id,
	}

	data, err := templateModel.Find(filter)

	if err != nil || len(data) == 0 {
		errorMap := map[string]string{
			"error": "error while deleting template",
		}
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	templateModel.DeleteById(data[0].ID.Hex())

	tempMap := make(map[string]interface{})

	tempMap["status"] = "success"

	json.NewEncoder(w).Encode(tempMap)
}
