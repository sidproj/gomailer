package controller

import (
	"fmt"
	"gomailer/models"
	"gomailer/utils"
	"net/http"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)
const viewsDir = "views"

var (
	loginHTML    = filepath.Join(viewsDir, "login.html")
	registerHTML = filepath.Join(viewsDir, "register.html")
)

type LoginData struct{
	Email string
	Password string
}

type RegisterData struct{
	FirstName string
	LastName string
	Email string
	Password string
	ConfirmPassword string
}

type AuthTemplateData struct{
	Error string
	CsrfToken string
}

func sendError(w http.ResponseWriter,template_path string) func(error string) {

	return func (error string){
		fmt.Printf("Error: %s!\n",error)
		e := AuthTemplateData{
				Error:error,
				CsrfToken:"",
			}
		token, err := utils.GenerateCSRFToken()
		if(err!=nil){
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		utils.SetCSRFCookie(w, token)
		utils.RenderTemplate(w,template_path,e)
	}
}

func LoginControllerGET(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GenerateCSRFToken()
	if(err!=nil){
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	utils.SetCSRFCookie(w, token)
	
	data := AuthTemplateData{
			Error:"",
			CsrfToken: token,
		}
	utils.RenderTemplate(w,loginHTML,data)
}

func LoginControllerPOST(w http.ResponseWriter, r *http.Request){
	lg := LoginData{Email: r.FormValue("Email"),Password: r.FormValue("Password")}

	sendLoginError := sendError(w,loginHTML)
	if lg.Email == "" || lg.Password == ""{
		sendLoginError("invalid request data")
		return
	}
	if !utils.VerifyCSRFToken(r) {
		sendLoginError("CSRF token mismatch")
		return
	}

	userModel,err := models.GetUserModel()
	
	if err != nil{
		fmt.Printf("Error while creating user model, error: %v",err )
		sendLoginError(err.Error())
		return		
	}

	filter := bson.M{"email":strings.ToLower(lg.Email)}

	users,err := userModel.Find(filter)

	if err != nil{
		sendLoginError(err.Error())
		return		
	}else if len(users) == 0{
		sendLoginError("no user found")
		return
	}else if !utils.CheckPasswordHash(lg.Password,users[0].Password){
		sendLoginError("incorrect password")
		return
	}

	// return login page with error if error
	fmt.Println("login success")

	jwtToken,err := utils.GenerateJWT(users[0].ID.Hex())
	if err != nil{
		sendLoginError("Error while creating jwt:"+err.Error())
		return
	}

	cookie := http.Cookie{
		Name:     "auth_jwt",
		Value:    jwtToken,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w,&cookie)
	http.Redirect(w,r,"/",http.StatusSeeOther)
}

func RegisterControllerGET(w http.ResponseWriter, r *http.Request){
	token ,err := utils.GenerateCSRFToken()
	if(err!=nil){
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	utils.SetCSRFCookie(w, token)
	e := AuthTemplateData{
			Error:"",
			CsrfToken: token,
		}
	utils.RenderTemplate(w,registerHTML,e)
}

func RegisterControllerPOST(w http.ResponseWriter,r *http.Request){
	rg := RegisterData{
		FirstName: r.FormValue("FirstName"),
		LastName: r.FormValue("LastName"),
		Email: strings.ToLower(r.FormValue("Email")),
		Password: r.FormValue("Password"),
		ConfirmPassword: r.FormValue("Confirm Password")}

	sendRegisterError := sendError(w,registerHTML)
	
	if rg.LastName=="" || rg.FirstName=="" || rg.Email == "" || rg.Password == "" || rg.ConfirmPassword == ""{
		fmt.Println(rg)
		sendRegisterError("Invalid request data")
		return
	}else if rg.Password != rg.ConfirmPassword {
		sendRegisterError("Password do not match")
		return
	}
	if !utils.VerifyCSRFToken(r) {
		sendRegisterError("CSRF token mismatch")
		return
	}

	userModel,err := models.GetUserModel()

	if err!=nil{
		sendRegisterError(err.Error())
		return	
	}

	hashedPassword,err := utils.HashPassword(rg.Password)
	
	if err!= nil{
		sendRegisterError(err.Error())	
		return
	}

	user := models.UserSchema{
		FirstName: strings.ToLower(rg.FirstName),
		LastName: rg.LastName,
		Email: rg.Email,
		Password: hashedPassword,
	}

	if err:=userModel.Save(&user);err!=nil{
		if(strings.Contains(err.Error(),"E11000")){
			sendRegisterError(fmt.Sprintf("Account already exists for email %s. Try different email!",user.Email))
			return
		}
		sendRegisterError(err.Error())
		return
	}

	http.Redirect(w,r,"/login",http.StatusSeeOther)
}