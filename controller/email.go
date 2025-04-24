package controller

import (
	"encoding/json"
	"fmt"
	"gomailer/models"
	"gomailer/service"
	"gomailer/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


var (
    authUserName = utils.GetEnvVariable("AWS_SMTP_USERNAME")
    authPassword = utils.GetEnvVariable("AWS_SMTP_PASSWORD")
	auth=service.SetupSMTPAuth(authUserName,authPassword,"email-smtp.eu-north-1.amazonaws.com")
)

type SendEmailData struct {
	TemplateName      string            `json:"template_name"`
	TemplateVariables map[string]string `json:"template_variables"`
}

func SendEmailControllerPOST(w http.ResponseWriter, r *http.Request) {

	// get user id from request header
	userID,err := primitive.ObjectIDFromHex(r.Header.Get("user_id"))
	if(err!=nil){
		errorMap := map[string]string{
			"error":"no user found!",
			"redirect":"/login",
		}

		json.NewEncoder(w).Encode(errorMap)
		return
	}

	// get request data
	decoder := json.NewDecoder(r.Body)
    var data SendEmailData
    err = decoder.Decode(&data)

    if err != nil {
        decode_err := map[string]string{
			"error":"invalid request",
			"description":"invalid request data or request data format",
		}
    	json.NewEncoder(w).Encode(decode_err)
		return
    }
	if( data.TemplateVariables==nil || len(data.TemplateName)==0){
		 decode_err := map[string]string{
			"error":"invalid request",
			"description":"invalid request data or request data format",
		}
    	json.NewEncoder(w).Encode(decode_err)
		return
	}

	// fetch template from template modal
	templateModel,err := models.GetTemplateModel()

	if(err!=nil){
		errorMap := map[string]string{
			"error":"error while getting template modal",
		}
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	filter := bson.M{"name":data.TemplateName,"user_id":userID}

	template,err := templateModel.Find(filter)

	if(err != nil){
		errorMap := map[string]string{
			"error":"error while finding template",
		}
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	fmt.Println(template)

	json.NewEncoder(w).Encode(template)
	return

	to := []string{"morisidhraj001@gmail.com"}
	err = service.SendMail(auth, to, "tesing", "Hello world!","sidharajdsa@gmail.com")
	if err != nil {
		fmt.Fprintf(w, "An error occured: %s", err.Error())
		return
	}
	fmt.Fprintf(w, "Email sent successfully!")
}