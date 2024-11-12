package controller

import (
	"fmt"
	"gomailer/utils"
	"html/template"
	"net/http"
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
	Email string
	Password string
	ConfirmPassword string
}

type ErrorData struct{
	Error string
}

func LoginControllerGet(w http.ResponseWriter, r *http.Request) {
	e := ErrorData{Error:""}
	t,_ := template.ParseFiles(login_html)
	t.Execute(w,e)
}

func LoginControllerPOST(w http.ResponseWriter, r *http.Request){
	lg := LoginData{Email: r.FormValue("Email"),Password: r.FormValue("Password")}

	fmt.Printf("email: %s password: %s\n",lg.Email,lg.Password)
	var error string  = ""

	if lg.Email == "" || lg.Password == ""{
		error = "Invalid request data!"
	}else if lg.Email != "sid@gmail.com" || lg.Password != "1234"{
		error = "Wrong email or password"
	}

	// return login page with error if error
	if(error != ""){
		fmt.Printf("Error: %s\n",error)
		e := ErrorData{Error:error}
		t,_ := template.ParseFiles(login_html)
		t.Execute(w,e)
		return
	}

	fmt.Println("login success")

	jwtToken,err := utils.GenerateJWT(lg.Email)
	if err != nil{
		error = "Error while creating jwt:"+err.Error()
	}

	// return login page with error if error
	if(error != ""){
		fmt.Printf("Error: %s\n",error)
		e := ErrorData{Error:error}
		t,_ := template.ParseFiles(login_html)
		t.Execute(w,e)
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
	rg := RegisterData{Email: r.FormValue("Email"),Password: r.FormValue("Password"),ConfirmPassword: r.FormValue("Confirm Password")}

	var error string = ""
	
	if rg.Email == "" || rg.Password == "" || rg.ConfirmPassword == ""{
		error = "invalid request data"
	}else if rg.Password != rg.ConfirmPassword {
		error = "password do not match"
	}

	if error != ""{
		fmt.Printf("Error: %s!\n",error)
		e := ErrorData{Error:error}
		t,_ := template.ParseFiles(register_html)
		t.Execute(w,e)
		return
	}

	http.Redirect(w,r,"/",http.StatusOK)

}