package repositories

import (
    "fmt"
    "gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/database"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
)

func InsertImage(image map[string]interface{}) interface{} {
    if len(image["ImageId"].(string)) == 0 || len(image["ImageDisplayName"].(string)) == 0 {
        return nil
    }

    insertFilter := bson.D{
        {Key: "image_id", Value: image["ImageId"]},
        {Key: "image_name", Value: image["ImageName"]},
        {Key: "image_display_name", Value: image["ImageDisplayName"]},
        {Key: "image_description", Value: image["ImageDescription"]},
        {Key: "published", Value: image["Published"]},
        {Key: "image_config", Value: image["ImageConfig"]},
        {Key: "image_read_root_password", Value: image["ImageReadRootPassword"]},
    }

    db, context, cancel := database.GetClient()
    collection := db.Database(database.DefaultDB).Collection(database.ImagesCollection)
    inserted, insertErr := collection.InsertOne(context, insertFilter)

    if insertErr != nil {
        log.Println("Error while inserting image to database!", insertErr)
        defer cancel()
        defer db.Disconnect(context)
        return nil
    }

    defer cancel()
    defer db.Disconnect(context)
    return inserted
}

func GetImages() []*database.Images {
    db, context, cancel := database.GetClient()
    collection := db.Database(database.DefaultDB).Collection(database.ImagesCollection)
    cursor, err := collection.Find(context, bson.D{})

    if err != nil {
        fmt.Println("Error while reading images!", err)
        defer cancel()
        defer db.Disconnect(context)
        return nil
    }

    var images []*database.Images
    for cursor.Next(context) {
        var elem database.Images
        err := cursor.Decode(&elem)
        if err != nil {
            fmt.Println("An error occurred while reading data!")
            defer cancel()
            defer cursor.Close(context)
            defer db.Disconnect(context)
            return nil
        }

        images = append(images, &elem)
    }

    defer cancel()
    defer db.Disconnect(context)
    return images
}

func GetImageByImageId(id interface{}) *database.Images {
    findFilter := bson.M{"image_id": id.(string)}

    var image database.Images

    db, context, cancel := database.GetClient()
    collection := db.Database(database.DefaultDB).Collection(database.ImagesCollection)
    err := collection.FindOne(context, findFilter).Decode(&image)

    if err != nil {
        fmt.Println("Error while reading image!", err)
        defer cancel()
        defer db.Disconnect(context)
        return nil
    }

    defer cancel()
    defer db.Disconnect(context)
    return &image
}

func GetImageById(id interface{}) *database.Images {
    documentId, _ := primitive.ObjectIDFromHex(id.(string))
    findFilter := bson.M{"_id": documentId}

    var image database.Images

    db, context, cancel := database.GetClient()
    collection := db.Database(database.DefaultDB).Collection(database.ImagesCollection)
    err := collection.FindOne(context, findFilter).Decode(&image)

    if err != nil {
        fmt.Println("Error while updating image!", err)
        defer cancel()
        defer db.Disconnect(context)
        return nil
    }

    defer cancel()
    defer db.Disconnect(context)
    return &image
}

func UpdateImageById(image map[string]interface{}) interface{} {
    documentId, _ := primitive.ObjectIDFromHex(image["Id"].(string))
    findFilter := bson.M{"_id": documentId}

    updateFilter := bson.M{"$set": bson.D{
        {Key: "image_id", Value: image["ImageId"]},
        {Key: "image_name", Value: image["ImageName"]},
        {Key: "image_display_name", Value: image["ImageDisplayName"]},
        {Key: "image_description", Value: image["ImageDescription"]},
        {Key: "published", Value: image["Published"]},
        {Key: "image_config", Value: image["ImageConfig"]},
        {Key: "image_read_root_password", Value: image["ImageReadRootPassword"]},
    }}

    db, context, cancel := database.GetClient()
    collection := db.Database(database.DefaultDB).Collection(database.ImagesCollection)
    _, err := collection.UpdateOne(context, findFilter, updateFilter)

    if err != nil {
        fmt.Println("Error while updating image!", err)
        defer cancel()
        defer db.Disconnect(context)
        return nil
    }

    defer cancel()
    defer db.Disconnect(context)
    return true
}

func DeleteImageById(id string) interface{} {
    documentId, _ := primitive.ObjectIDFromHex(id)
    findFilter := bson.M{"_id": documentId}

    db, context, cancel := database.GetClient()
    collection := db.Database(database.DefaultDB).Collection(database.ImagesCollection)
    _, err := collection.DeleteOne(context, findFilter)

    if err != nil {
        fmt.Println("Error while updating image!", err)
        defer cancel()
        defer db.Disconnect(context)
        return nil
    }

    defer cancel()
    defer db.Disconnect(context)
    return true
}

func GetPublishedImages() []*database.Images {
    opts := options.Find().SetProjection(bson.D{

        {
            Key:   "_id",
            Value: 0,
        },
        {
            Key:   "image_name",
            Value: 0,
        },
        {
            Key:   "image_description",
            Value: 0,
        },
        {
            Key:   "image_config",
            Value: 0,
        },
    })

    db, context, cancel := database.GetClient()
    collection := db.Database(database.DefaultDB).Collection(database.ImagesCollection)
    cursor, err := collection.Find(context, bson.D{{
        Key:   "published",
        Value: "true",
    }}, opts)

    if err != nil {
        fmt.Println("Error while reading images!", err)
        defer cancel()
        defer db.Disconnect(context)
        return nil
    }

    var images []*database.Images
    for cursor.Next(context) {
        var elem database.Images
        err := cursor.Decode(&elem)
        if err != nil {
            fmt.Println("An error occurred while reading data!")
            defer cancel()
            defer cursor.Close(context)
            defer db.Disconnect(context)
            return nil
        }

        images = append(images, &elem)
    }

    defer cancel()
    defer db.Disconnect(context)
    return images
}