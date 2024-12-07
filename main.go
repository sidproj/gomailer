package main

import (
	"fmt"
	"gomailer/controller"
	"gomailer/mango"
	"gomailer/middleware"
	"gomailer/models"
	"gomailer/router"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func main(){
	addingRoutes()
	router.LoadRoutes()

	mangoClient := mango.MongoConnect("mongodb://localhost:27017/","gomailer")
	defer mangoClient.CloseConn()
	
	UserModel := mango.CreateModel[models.UserSchema]("users")

	user,err := UserModel.FindById("674b78e0d14defb43b0dcc01")

	if err !=nil{
		if err == mongo.ErrNoDocuments{
			fmt.Println("No documents found with this given id.")
		}else{
			panic(err)
		}
	}

	user.Email = "Hello@gmail.com"

	// UserModel.Save(user)

	users,err := UserModel.Find(bson.M{"$or":bson.A{
		bson.M{"email":"test@gmail.com"},
		bson.M{"first_name":"john"},
	}})

	if err != nil{
		fmt.Println(err)
	}

	for _,u := range users{
		fmt.Println(u)
	}

	// fmt.Println("Server is running at http://localhost:8080")
	// log.Fatal(http.ListenAndServe(":8080",nil))

}