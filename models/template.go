package models

import (
	"gomailer/mango"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TemplateSchema struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id,omitempty"`
	TemplateContent string `bson:"template_content"`
	TemplateVariables []string `bson:"template_variables,omitempty"`
}

var templateModel *mango.GenericCollectionModel[TemplateSchema]

func GetTemplateModel()(*mango.GenericCollectionModel[TemplateSchema],error){
	if(templateModel != nil) {
		return templateModel,nil
	}

	model := mango.CreateModel[TemplateSchema]("templates")
	templateModel =&model
	return &model,nil
}
