package controller

import (
	"fmt"
	"gomailer/models"
	"gomailer/utils"
	"html/template"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

const(
	login_html = "views\\login.html"
	register_html = "views\\register.html"
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

type ErrorData struct{
	Error string
}

func sendError(w http.ResponseWriter,template_path string) func(error string) {

	return func (error string){
		fmt.Printf("Error: %s!\n",error)
		e := ErrorData{Error:error}
		t,_ := template.ParseFiles(template_path)
		t.Execute(w,e)
	}
}

func LoginControllerGet(w http.ResponseWriter, r *http.Request) {
	e := ErrorData{Error:""}
	t,_ := template.ParseFiles(login_html)
	t.Execute(w,e)
}

func LoginControllerPOST(w http.ResponseWriter, r *http.Request){
	lg := LoginData{Email: r.FormValue("Email"),Password: r.FormValue("Password")}

	sendLoginError := sendError(w,login_html)
	if lg.Email == "" || lg.Password == ""{
		sendLoginError("invalid request data")
		return
	}
	userModel,err := models.GetUserModel()
	
	if err != nil{
		fmt.Printf("Error while creating user model, error: %v",err )
		sendLoginError(err.Error())
		return		
	}

	filter := bson.M{"email":lg.Email}

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

	jwtToken,err := utils.GenerateJWT(lg.Email)
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

func RegisterControllerGet(w http.ResponseWriter, r *http.Request){
	e := ErrorData{Error:""}
	t,_ := template.ParseFiles(register_html)
	t.Execute(w,e)
}

func RegisterControllerPost(w http.ResponseWriter,r *http.Request){
	rg := RegisterData{
		FirstName: r.FormValue("FirstName"),
		LastName: r.FormValue("LastName"),
		Email: r.FormValue("Email"),
		Password: r.FormValue("Password"),
		ConfirmPassword: r.FormValue("Confirm Password")}

	sendRegisterError := sendError(w,register_html)
	
	if rg.LastName=="" || rg.FirstName=="" || rg.Email == "" || rg.Password == "" || rg.ConfirmPassword == ""{
		sendRegisterError("invalid request data")
	}else if rg.Password != rg.ConfirmPassword {
		sendRegisterError("password do not match")
	}

	
	hashedPassword,err := utils.HashPassword(rg.Password)
	
	if err!= nil{
		sendRegisterError(err.Error())	
		return
	}

	userModel,err := models.GetUserModel()

	if err!=nil{
		sendRegisterError(err.Error())
		return	
	}

	user := models.UserSchema{
		FirstName: rg.FirstName,
		LastName: rg.LastName,
		Email: rg.Email,
		Password: hashedPassword,
	}

	if err:=userModel.Save(user);err!=nil{
		sendRegisterError(err.Error())
		return
	}

	http.Redirect(w,r,"/",http.StatusSeeOther)

}