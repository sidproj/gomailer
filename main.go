package main

import (
	"fmt"
	"gomailer/controller"
	"gomailer/mango"
	"gomailer/middleware"
	"gomailer/models"
	"gomailer/router"
	"log"
	"net/http"
)

/*
var (
    authUserName = goDotEnvVariable("AWS_SMTP_USERNAME")
    authPassword = goDotEnvVariable("AWS_SMTP_PASSWORD")
	auth=service.SetupSMTPAuth(authUserName,authPassword,"email-smtp.eu-north-1.amazonaws.com")
)

func handler1(w http.ResponseWriter,r *http.Request){
	to := []string{"morisidhraj001@gmail.com"}
	err := service.SendMail(auth,to,"tesing","Hello world!")
	if(err != nil){
		fmt.Fprintf(w,"An error occured: %s",err.Error())
		return
	}
	fmt.Fprintf(w,"Email sent successfully!")
}
*/


func addingRoutes(){
	// router.Get("/",handler1)
	router.Get("/",controller.HomeController)
	router.Get("/login",controller.LoginControllerGet)
	router.Post("/login",controller.LoginControllerPOST)

	router.Get("/register",controller.RegisterControllerGet)
	router.Post("/register",controller.RegisterControllerPost)

	router.Get("/secret",middleware.AuthMiddlewareUser(controller.HomeController))
}

func loadModels(){
	_,err := models.GetUserModel()
	
	if err!=nil{
		fmt.Printf("Error while creating user model. Error: %v",err)
	}
	fmt.Println("Loaded user model")
}

func main(){
	addingRoutes()
	router.LoadRoutes()

	mangoClient := mango.MongoConnect("mongodb://localhost:27017/","gomailer")
	defer mangoClient.CloseConn()

	// load models
	loadModels()

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080",nil))

}