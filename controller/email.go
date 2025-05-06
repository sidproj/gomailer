package controller

import (
	"encoding/json"
	"fmt"
	"gomailer/models"
	"gomailer/service"
	"gomailer/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)


var (
    authUserName = utils.GetEnvVariable("AWS_SMTP_USERNAME")
    authPassword = utils.GetEnvVariable("AWS_SMTP_PASSWORD")
	auth=service.SetupSMTPAuth(authUserName,authPassword,"email-smtp.eu-north-1.amazonaws.com")
)

type SendEmailData struct {
	Email 			  string 			`json:"email"`
	Password 		  string 			`json:"password"`
	TemplateName      string            `json:"template_name"`
	TemplateVariables map[string]string `json:"template_variables"`
}

func sendJsonError(w http.ResponseWriter,error string,description string) {
	fmt.Printf("Error: %s!\n",error)
	errorMap := map[string]string{
		"error":error,
		"description":description,
	}
	json.NewEncoder(w).Encode(errorMap)
}

func validateMapKeys(keys []string,m2 map[string]string)bool{
	fmt.Println(keys)
	for _,k:= range keys{
		fmt.Println(k,m2[k])
		_,exists := m2[k]
		if(!exists) {
			return false
		}
	}
	return true
}

func SendEmailControllerPOST(w http.ResponseWriter, r *http.Request) {

	// get request data
	decoder := json.NewDecoder(r.Body)
    var data SendEmailData
    err := decoder.Decode(&data)

    if ( err != nil || 
		data.TemplateVariables == nil || 
		len(data.TemplateName) == 0 || 
		len(data.Email) == 0 || 
		len(data.Password) == 0 ) {
		sendJsonError(w,"invalid request","invalid request data or request data format")
        return
    }

	// get email and password from request json
	userModel,err:= models.GetUserModel()

	if(err!=nil){
		sendJsonError(w,"user model creation error",err.Error())
		return
	}

	filter := bson.M{"email":data.Email}

	users,err := userModel.Find(filter)

	if err != nil{
		sendJsonError(w,"login error",err.Error())
		return		
	}else if len(users) == 0{
		sendJsonError(w,"login error","no user found")
		return
	}else if !utils.CheckPasswordHash(data.Password,users[0].Password){
		sendJsonError(w,"login error","incorrect password")
		return
	}
	

	if( data.TemplateVariables==nil || len(data.TemplateName)==0){
		sendJsonError(w,"invalid request","invalid request data or request data format")
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

	filter = bson.M{"name":data.TemplateName,"user_id":users[0].ID}

	template,err := templateModel.Find(filter)

	if(err != nil){
		errorMap := map[string]string{
			"error":"error while finding template",
		}
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	if(!validateMapKeys(template[0].TemplateVariables,data.TemplateVariables)){
		errorMap := map[string]string{
			"error":"invalid template variables",
		}
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	replaceContent := utils.ReplaceTemplateVariables(template[0].TemplateContent,data.TemplateVariables)

	to := []string{"morisidhraj001@gmail.com"}
	err = service.SendMail(auth, to, "tesing replacement", replaceContent,"sidharajdsa@gmail.com")
	if err != nil {
		errorMap := map[string]string{
			"error":fmt.Sprintf("An error occured: %s",err.Error()),
		}
		json.NewEncoder(w).Encode(errorMap)
		return
	}
	msgMap := map[string]string{
		"status":"Success",
		"msg":"Email sent successfully!",
	}
	json.NewEncoder(w).Encode(msgMap)
}