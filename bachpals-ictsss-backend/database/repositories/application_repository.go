package repositories

import (
    "fmt"
    "github.com/spf13/viper"
    "gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/database"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "os"
)

func InitializeDefaultAdministrator() interface{} {
    db, context, cancel := database.GetClient()
    collection := db.Database(database.DefaultDB).Collection(database.AdminCollection)

    filter := bson.D{{"user_id", viper.GetString("IKT_STACK_DEFAULT_ADMIN")}}

    var defaultAdmin database.Admin
    if err := collection.FindOne(context, filter).Decode(&defaultAdmin); err != nil {
        if err == mongo.ErrNoDocuments {
            defer cancel()
            defer db.Disconnect(context)
            fmt.Printf("Default administrator doesn't exist! Adding one provided in the config\n")

            insertData := bson.D{{Key: "user_id", Value: viper.GetString("IKT_STACK_DEFAULT_ADMIN")}, {Key: "name", Value: "Default admin"}}

            collection := db.Database(database.DefaultDB).Collection(database.AdminCollection)
            inserted, insertErr := collection.InsertOne(context, insertData)

            if insertErr != nil {
                defer cancel()
                defer db.Disconnect(context)
                fmt.Printf("Error while adding default administrator!\n %+v\n", insertErr)
                os.Exit(0)
            }

            if inserted.InsertedID != nil {
                defer cancel()
                defer db.Disconnect(context)

                collection := db.Database(database.DefaultDB).Collection(database.ServerCollection)
                filter := bson.D{{"first_run", 1}}
                inserted, err := collection.InsertOne(context, filter)

                if err != nil {
                    defer cancel()
                    defer db.Disconnect(context)
                    fmt.Printf("Failed updating initialization state!\n %+v\n", err)
                    return nil
                }

                if inserted.InsertedID != nil {
                    defer cancel()
                    defer db.Disconnect(context)
                    fmt.Printf("Added default administrator successfully!\n %s\n", viper.GetString("IKT_STACK_DEFAULT_ADMIN"))
                    return nil
                }
            }
        }

        defer cancel()
        defer db.Disconnect(context)
        fmt.Printf("Error while initializing default administrator!\n %+v\n", err)
        os.Exit(0)
    }

    if defaultAdmin.UserId != "" {
        fmt.Printf("Application has already been initialized!\n")
    }

    return nil
}

// True if server is already initialized, false if not.
func CheckInitializationStatus() bool {
    db, context, cancel := database.GetClient()
    collection := db.Database(database.DefaultDB).Collection(database.ServerCollection)

    filter := bson.D{{Key: "first_run", Value: 1}}

    var server database.Application
    err := collection.FindOne(context, filter).Decode(&server)

    if err != nil {
        defer cancel()
        defer db.Disconnect(context)

        if err == mongo.ErrNoDocuments {
            fmt.Printf("No initialization status document in database! Creating document!\n %+v\n", err)

            inserted, insertErr := collection.InsertOne(context, bson.D{{Key: "first_run", Value: 1}})

            if insertErr != nil {
                fmt.Printf("Error while creating initialization document!\n %+v\n", insertErr)
                defer cancel()
                defer db.Disconnect(context)
                os.Exit(0)
            }

            if inserted.InsertedID != "" {
                defer cancel()
                defer db.Disconnect(context)

                fmt.Printf("Successfully created initialization document!\n %s\n", inserted.InsertedID)
                return true
            }
        } else {
            defer cancel()
            defer db.Disconnect(context)
            fmt.Printf("Error while checking server initialization status!\n", err)
            os.Exit(0)
        }
    }

    if server.FirstRun != 1 {
        defer cancel()
        defer db.Disconnect(context)
        return false
    }
    return true
}

func RemoveInitializationStatus() {
    db, context, cancel := database.GetClient()
    if err := db.Database(database.DefaultDB).Collection(database.ServerCollection).Drop(context); err != nil {
        fmt.Printf("Could not drop sever collection!\n")
        defer cancel()
        defer db.Disconnect(context)
        os.Exit(0)
    }

    if err := db.Database(database.DefaultDB).Collection(database.AdminCollection).Drop(context); err != nil {
        fmt.Printf("Could not drop administrators collection!\n")
        defer cancel()
        defer db.Disconnect(context)
        os.Exit(0)
    }

    if err := db.Database(database.DefaultDB).Collection(database.VmCollection).Drop(context); err != nil {
        fmt.Printf("Could not drop vm collection!\n")
        defer cancel()
        defer db.Disconnect(context)
        os.Exit(0)
    }

    defer cancel()
    defer db.Disconnect(context)

    fmt.Printf("System has been reinitialized!\n")
    os.Exit(0)
}