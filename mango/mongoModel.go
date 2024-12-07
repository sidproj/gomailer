package mango

import (
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GenericCollectionModel[T any] struct {
	collectionName  string
	mongoCollection *mongo.Collection
	// context context.Context
}

// Save function recives any interface type parameter and
// checks if the parameter is sturct, returns error if not.
// If the struct has ID filed populated then it updates document with the struct.
// Else inserts a new document into the collection and set's its ID in the struct passed.
func (collection *GenericCollectionModel[T])Save(model interface{})error{

	reflectedModel := reflect.ValueOf(model)
    
    if reflectedModel.Kind() == reflect.Ptr{
        reflectedModel = reflectedModel.Elem()
    }

    if reflectedModel.Kind() != reflect.Struct{
        return errors.New("model schema must be a struct")
    }

    idField := reflectedModel.FieldByName("ID") 

    if !idField.IsValid(){
        return errors.New("model schema does not have an 'ID' field")
    }

    if idField.IsZero(){
        
        result,err := collection.mongoCollection.InsertOne(dbClient.context,model)

        if err!= nil {
            return err
        }

        // update id in struct
        if idField.CanSet() && result.InsertedID != nil{
            idField.Set(reflect.ValueOf(result.InsertedID))
        }

        // remove later
        fmt.Printf("Inserted data in database for collection: %s with id: %v\n",collection.collectionName,result.InsertedID)

    }else{
 
        filter := bson.M{"_id":idField.Interface()}
        update := bson.M{"$set":model}

        opts := options.Update().SetUpsert(true)

        _,err := collection.mongoCollection.UpdateOne(dbClient.context,filter,update,opts)

        if err != nil{
            return err
        }

        // remove later
        fmt.Printf("Updated data in database for collection: %s with id: %v\n",collection.collectionName,idField.Interface())
    }
    return nil
}

func (collection *GenericCollectionModel[T])Find(filter bson.M)([]T,error){

	var data []T

	cursor,err := collection.mongoCollection.Find(dbClient.context,filter)
	if err != nil{
		return nil,err
	}

	if err := cursor.All(dbClient.context,&data); err != nil{
		return nil,err
	}

	return data,nil
}

// FindById returns model(struct) by taking id as a parameter. It checks if
// id is a valid ObjectID. Returns model(struct) and err.
func (collection *GenericCollectionModel[T])FindById(id string)(T,error){

	var model T
    // Convert the string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model, errors.New("invalid ObjectID")
	}

    filter := bson.M{"_id":objectID}

    if err := collection.mongoCollection.FindOne(dbClient.context,filter).Decode(&model); err != nil{
        return model,err
    }
    
    fmt.Println("Func data: ",model)

    return model,nil
}