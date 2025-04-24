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

func GetUserModel()(*mango.GenericCollectionModel[UserSchema],error){

	if(userModel != nil) {
		return userModel,nil
	}

	model := mango.CreateModel[UserSchema]("users")
	if err:=model.CreateIndex([]string{"email"});err!=nil{
		return nil,err
	}
	userModel = &model
	return &model,nil
}

