package database

import (
    "context"
    "github.com/spf13/viper"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
    "time"
)

const DefaultDB = "ikt-stack"
const VmCollection = "virtual_machines"
const AdminCollection = "administrators"
const ServerCollection = "server"
const ImagesCollection = "images"

type MongoHandler struct {
    Mongo *mongo.Collection
}

func GetClient() (*mongo.Client, context.Context, context.CancelFunc) {
    client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("IKT_STACK_DB_URL")))
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        defer cancel()
        log.Fatal(err)
    }

    return client, ctx, cancel
}