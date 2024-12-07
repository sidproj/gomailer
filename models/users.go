package models

import (
	"gomailer/mango"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSchema struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string `bson:"first_name"`
	LastName string `bson:"last_name"`
	Email string `bson:"email"`
	Password string `bson:"password"`
}

var userModel *mango.GenericCollectionModel[UserSchema]

func GetUserModel()*mango.GenericCollectionModel[UserSchema]{

	if(userModel != nil) {
		return userModel
	}

	model := mango.CreateModel[UserSchema]("users")
	userModel = &model
	return &model
}

