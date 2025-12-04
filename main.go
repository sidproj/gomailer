package main

import (
	"fmt"
	"gomailer/controller"
	"gomailer/mango"
	"gomailer/middleware"
	"gomailer/models"
	"gomailer/utils"
	"log"
	"net/http"

	router "github.com/sidproj/grouter"
)

func addingRoutes() {
	router.Get("/", controller.HomeController)
	router.Get("/login", controller.LoginControllerGET)
	router.Post("/login", controller.LoginControllerPOST)

	router.Get("/register", controller.RegisterControllerGET)
	router.Post("/register", controller.RegisterControllerPOST)

	// template testing
	router.Get("/template", middleware.AuthMiddlewareUser(controller.TemplateControllerGET))
	router.Get("/template/create", middleware.AuthMiddlewareUser(controller.CreateTemplateControllerGET))
	router.Post("/template/create", middleware.AuthMiddlewareUser(controller.CreateTemplateControllerPOST))
	router.Get("/template/edit", middleware.AuthMiddlewareUser(controller.EditTemplateControllerGET))
	router.Post("/template/edit", middleware.AuthMiddlewareUser(controller.EditTemplateControllerPOST))
	router.Delete("/template/delete", middleware.AuthMiddlewareUser(controller.DeleteTemplateControllerDELETE))

	// static
	router.Get("/about", controller.AboutController)
	router.Get("/privacy-policy", controller.PrivacyPollicyController)

	// this endpoint will be available to other websites
	router.Post("/sendmail", middleware.PublicRouteMiddleware(controller.SendEmailControllerPOST))

}

func loadModels() {
	_, err := models.GetUserModel()

	if err != nil {
		fmt.Printf("Error while creating user model. Error: %v", err)
	}
	fmt.Println("Loaded user model")
}

func main() {
	addingRoutes()
	router.LoadRoutes()

	mangoClient := mango.MongoConnect(utils.GetEnvVariable("MONGODB_ATLAS_URI"), "gomailer")
	defer mangoClient.CloseConn()

	// load models
	loadModels()

	fmt.Println("Server is running at http://localhost:8080")
	// Serve static assets like images, CSS, JS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
