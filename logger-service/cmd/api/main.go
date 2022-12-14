package main

import (
	"context"
	"fmt"
	"github.com/Christomesh/go-microservice/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	//connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panicln()
	}
	client = mongoClient

	// create a context in order to disconnect

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}
	//go app.serve()
	srv := http.Server{
		Addr:    fmt.Sprintf("%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panicln(err)
	}

}

//func (app *Config) serve() {
//	srv := http.Server{
//		Addr:    fmt.Sprintf("%s", webPort),
//		Handler: app.routes(),
//	}
//
//	err := srv.ListenAndServe()
//	if err != nil {
//		log.Panicln(err)
//	}
//}

func connectToMongo() (*mongo.Client, error) {
	// create connect options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting", err)
		return nil, err
	}
	return c, nil

}
