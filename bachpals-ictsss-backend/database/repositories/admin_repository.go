package repositories

import (
    "fmt"
    "gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/database"
    "go.mongodb.org/mongo-driver/bson"
    "log"
)

func InsertAdmin(userId string, name string) interface{} {
    if len(userId) == 0 {
        return nil
    }

    if len(name) == 0 {
        return nil
    }

    insertData := bson.D{{Key: "user_id", Value: userId}, {Key: "name", Value: name}}

    db, context, cancel := database.GetClient()
    collection := db.Database(database.DefaultDB).Collection(database.AdminCollection)
    inserted, err := collection.InsertOne(context, insertData)

    if err != nil {
        defer db.Disconnect(context)
        defer cancel()
        return nil
    }

    defer db.Disconnect(context)
    defer cancel()

    return inserted
}

func ReadAdmins() []*database.Admin {
    db, context, cancel := database.GetClient()
    collection := db.Database(database.DefaultDB).Collection(database.AdminCollection)

    filter := bson.D{}

    var userData []*database.Admin
    cursor, err := collection.Find(context, filter)

    if err != nil {
        defer cancel()
        defer db.Disconnect(context)
        return nil
    }

    for cursor.Next(context) {
        var elem database.Admin
        err := cursor.Decode(&elem)
        if err != nil {
            fmt.Println("An error occurred while reading data!")
            defer cancel()
            defer db.Disconnect(context)
            defer cursor.Close(context)
            return nil
        }

        userData = append(userData, &elem)
    }

    defer cancel()
    defer db.Disconnect(context)
    defer cursor.Close(context)

    return userData
}

func ReadAdminById(userId string) interface{} {
    if len(userId) == 0 {
        return nil
    }

    db, context, cancel := database.GetClient()
    collection := db.Database(database.DefaultDB).Collection(database.AdminCollection)

    filter := bson.M{"user_id": bson.M{"$eq": userId}}

    var userData database.Admin
    err := collection.FindOne(context, filter).Decode(&userData)

    if err != nil {
        defer cancel()
        defer db.Disconnect(context)
        return nil
    }

    defer cancel()
    defer db.Disconnect(context)

    return userData
}

func DeleteAdminById(userId string) interface{} {
    if len(userId) == 0 {
        return nil
    }

    db, context, cancel := database.GetClient()
    collection := db.Database(database.DefaultDB).Collection(database.AdminCollection)

    filter := bson.M{"user_id": userId}
    res, err := collection.DeleteOne(context, filter)

    if err != nil {
        defer cancel()
        defer db.Disconnect(context)
        return nil
    }

    defer cancel()
    defer db.Disconnect(context)

    return int(res.DeletedCount)
}

func UpdateAdminById(id string, newId string, name string) bool {
    if len(id) == 0 {
        return false
    }

    log.Println(bson.M{"$set": bson.M{"user_id": newId, "name": name}})

    findFilter := bson.M{"user_id": id}
    updateFilter := bson.M{"$set": bson.M{"user_id": newId, "name": name}}

    db, context, cancel := database.GetClient()
    vms := db.Database(database.DefaultDB).Collection(database.AdminCollection)
    _, err := vms.UpdateOne(context, findFilter, updateFilter)

    if err != nil {
        log.Println("Result: ", err)
        defer cancel()
        defer db.Disconnect(context)
        return false
    }

    defer cancel()
    defer db.Disconnect(context)

    return true
}