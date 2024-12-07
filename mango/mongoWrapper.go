package mango

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDbClient struct{
    client *mongo.Client
    database string
    context context.Context
    Cancel context.CancelFunc
}

// implementing singolton pattern
var dbClient *MongoDbClient

// make this function more robust later
func MongoConnect (uri string,database string)*MongoDbClient{

    if dbClient !=nil {
        return dbClient
    }

    // ctx will be used to set deadline for process, here 
    // deadline will of 30 seconds.
    ctx, cancel := context.WithTimeout(context.Background(), 
                                       30 * time.Second)
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

     
    if err != nil{
        panic(err)
    }

    dbClient = &MongoDbClient{
        client: client,
        context: ctx,
        Cancel: cancel,
        database: database,
    }

    dbClient.pingMongo()

    return dbClient
}

// make this function more robust later
func (dbClient *MongoDbClient)CloseConn(){
    
    dbClient.checkConnection()
    
    defer dbClient.Cancel()
    defer func(){
        if err := dbClient.client.Disconnect(dbClient.context);err != nil{
            panic(err)
        }
    }()
}

// utility function
func (dbClient *MongoDbClient)checkConnection(){
    if dbClient == nil || dbClient.client == nil {
        panic("No available connection!")
    }
}

// utility function
func (dbClient *MongoDbClient)pingMongo()error{
    if err := dbClient.client.Ping(dbClient.context,readpref.Primary()); err != nil{
     return err
    }
    fmt.Println("Database connection successful.")
    return nil
}

func (dbClient *MongoDbClient)getDatabase()*mongo.Database{
    dbClient.checkConnection()
    return dbClient.client.Database(dbClient.database)
}

// model functions: Insert of update the values of a collection provided in the first and data provided in the second parameter
func (dbClient *MongoDbClient)Save(collection string,model interface{})error{

    dbClient.checkConnection()

    reflectedModel := reflect.ValueOf(model)
    
    if reflectedModel.Kind() == reflect.Ptr{
        reflectedModel = reflectedModel.Elem()
    }

    if reflectedModel.Kind() != reflect.Struct{
        return errors.New("Model schema must be a struct")
    }

    idField := reflectedModel.FieldByName("ID") 

    if !idField.IsValid(){
        return errors.New("Model schema does not have an 'ID' field")
    }

    mongoCollection := dbClient.getDatabase().Collection(collection)

    if idField.IsZero(){
        
        result,err := mongoCollection.InsertOne(dbClient.context,model)

        if err!= nil {
            return err
        }

        // update id in struct
        if idField.CanSet() && result.InsertedID != nil{
            idField.Set(reflect.ValueOf(result.InsertedID))
        }

        // remove later
        fmt.Printf("Inserted data in database for collection: %s with id: %v\n",collection,result.InsertedID)

    }else{
 
        filter := bson.M{"_id":idField.Interface()}
        update := bson.M{"$set":model}

        opts := options.Update().SetUpsert(true)

        _,err := mongoCollection.UpdateOne(dbClient.context,filter,update,opts)

        if err != nil{
            return err
        }

        // remove later
        fmt.Printf("Updated data in database for collection: %s with id: %v\n",collection,idField.Interface())
    }
    return nil
}

// to get ORM like mongoose's model function
func CreateModel[T any](collection string)GenericCollectionModel[T]{

    dbClient.checkConnection()

    modelCollection := dbClient.getDatabase().Collection(collection)

    return GenericCollectionModel[T]{
        collectionName: collection,
        mongoCollection: modelCollection,
    }

}