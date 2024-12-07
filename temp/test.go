package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func connect (uri string)(*mongo.Client, context.Context, context.CancelFunc, error){
    
    // ctx will be used to set deadline for process, here 
    // deadline will of 30 seconds.
    ctx, cancel := context.WithTimeout(context.Background(), 
                                       30 * time.Second)
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    return client, ctx, cancel, err

}

func closeConn(client * mongo.Client,ctx context.Context,cancel context.CancelFunc){
    defer cancel()
    defer func(){
        if err := client.Disconnect(ctx);err != nil{
            panic(err)
        }
    }()
}

func pingMongo(client * mongo.Client, ctx context.Context)error{
    if err := client.Ping(ctx,readpref.Primary()); err != nil{
     return err
    }
    fmt.Println("Connected Successfully")
    return nil
}
    
func SetupMongodbConn(){
	// mongodb_uri := GetEnvVariable("MONGODB_URI")
    mongodb_uri := "mongodb://localhost:27017/"
	client, context,cancel, err := connect(mongodb_uri)
    if err != nil{
        panic(err)
    }

    defer closeConn(client,context,cancel)

    pingMongo(client,context)

    usersCollection := client.Database("gomailer").Collection("users")
    
    cursor,err := usersCollection.Find(context,bson.D{})

    if err != nil{
        panic(err)
    }

    var results []bson.D

    if err:=cursor.All(context,&results);err != nil{
        panic(err)
    }

    for _,doc:=range results{
		data := make(map[string]any)

		marshaledData,err := bson.Marshal(doc)
		if err != nil{
			panic(err)
		}

		if err := bson.Unmarshal(marshaledData,data);err!= nil{
			panic(err)
		}
        fmt.Println("Name: ",data["name"])
    }

}

func main(){
	SetupMongodbConn()
}