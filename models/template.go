package models

import (
	"errors"
	"gomailer/mango"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TemplateSchema struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id,omitempty"`
	Name string `bson:"name"`
	TemplateContent string `bson:"template_content"`
	TemplateVariables []string `bson:"template_variables,omitempty"`
}

func (t TemplateSchema) Validate() error {
	if t.Name == "" {
		return errors.New("name field is required and cannot be empty")
	}
	if t.UserID.String() == ""{
		return errors.New("userid field is required and cannot be empty")	
	}
	if(t.TemplateContent == ""){
		return errors.New("templatecontent field is required and cannot be empty")
	}
	return nil
}

var templateModel *mango.GenericCollectionModel[TemplateSchema]

func GetTemplateModel()(*mango.GenericCollectionModel[TemplateSchema],error){
	if(templateModel != nil) {
		return templateModel,nil
	}

	model := mango.CreateModel[TemplateSchema]("templates")
	model.CreateIndex([]string{"user_id","name"})
	templateModel =&model
	return &model,nil
}
