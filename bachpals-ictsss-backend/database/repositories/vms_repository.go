package repositories

import (
    "fmt"
    "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
    "gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/database"
    "gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/gopher"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
    "strings"
    "time"
)

func InsertMultipleVms(server *servers.Server, serverIp string, serverName string, users string, serverImage string) *mongo.InsertManyResult {

    var documents []interface{}

    time := time.Now()
    for _, k := range strings.Split(users, ",") {
        documents = append(documents, bson.D{
            {Key: "server_ip", Value: serverIp},
            {Key: "server_image", Value: serverImage},
            {Key: "server_name", Value: serverName},
            {Key: "user_id", Value: k},
            {Key: "server_status", Value: database.VirtualMachineStatusActive},
            {Key: "server_id", Value: server.ID},
            {Key: "created", Value: time},
        })
    }

    db, context, cancel := database.GetClient()
    vms := db.Database(database.DefaultDB).Collection(database.VmCollection)
    insertResponse, insertError := vms.InsertMany(context, documents)

    if insertError != nil {
        client := gopher.GetClient()

        // If for any reason a server could not be stored in the database, remove it from bare metal.
        result := servers.Delete(client, server.ID)

        // Do we want to print a response or fail silently?
        // This should probably be logged into some kind of logging system.
        if result.ExtractErr() != nil {
            fmt.Println("This error occurred while deleting vm,"+
                "after database failed insertion of a new vm.", result.ExtractErr())
        }
        defer cancel()
        defer db.Disconnect(context)
        return nil
    }
    defer cancel()
    defer db.Disconnect(context)

    return insertResponse
}

func InsertVm(server *servers.Server, serverIp string, serverName string, userId string) *mongo.InsertOneResult {

    virtualMachineMetadata := bson.D{
        {Key: "server_ip", Value: serverIp},
        {Key: "server_image", Value: server.Image},
        {Key: "server_name", Value: serverName},
        {Key: "user_id", Value: userId},
        {Key: "server_status", Value: database.VirtualMachineStatusActive},
        {Key: "server_id", Value: server.ID},
        {Key: "created", Value: time.Now()},
    }

    db, context, cancel := database.GetClient()
    vms := db.Database(database.DefaultDB).Collection(database.VmCollection)
    insertResponse, insertError := vms.InsertOne(context, virtualMachineMetadata)

    if insertError != nil {
        client := gopher.GetClient()

        // If for any reason a server could not be stored in the database, remove it from bare metal.
        result := servers.Delete(client, server.ID)

        // Do we want to print a response or fail silently?
        // This should probably be logged into some kind of logging system.
        if result.ExtractErr() != nil {
            fmt.Println("This error occurred while deleting vm,"+
                "after database failed insertion of a new vm.", result.ExtractErr())
        }
        defer cancel()
        defer db.Disconnect(context)
        return nil
    }
    defer cancel()
    defer db.Disconnect(context)

    return insertResponse
}

func GetVMS() []database.VirtualMachine {
    groupStage := []bson.M{{
        "$group": bson.M{
            "_id": "$server_id",
            "doc": bson.M{
                "$first": "$$ROOT",
            },
        },
    }}

    db, context, cancel := database.GetClient()
    vms := db.Database(database.DefaultDB).Collection(database.VmCollection)
    cur, err := vms.Aggregate(context, groupStage)

    if err != nil {
        defer cancel()
        defer db.Disconnect(context)
        defer cur.Close(context)
        return nil
    }

    var tmp []map[string]interface{}
    err = cur.All(context, &tmp)
    if err != nil {
        defer cancel()
        defer db.Disconnect(context)
        defer cur.Close(context)
        return nil
    }

    var virtualMachines []database.VirtualMachine
    for _, v := range tmp {
        var t database.VirtualMachine
        d, _ := bson.Marshal(v["doc"])
        bson.Unmarshal(d, &t)
        virtualMachines = append(virtualMachines, t)
    }

    var data []database.VirtualMachine
    for _, v := range virtualMachines {
        members := GetVmGroupMembers(v.ServerId)
        v.GroupMembers = members

        image := GetImageByImageId(v.ServerImage)
        v.ImageReadRootPassword = image.ImageReadRootPassword
        v.ImageDisplayName = image.ImageDisplayName

        data = append(data, v)
    }

    defer cancel()
    defer db.Disconnect(context)
    defer cur.Close(context)
    return data
}

func GetVMByUserId(id interface{}) interface{} {
    var virtualMachines []interface{}
    findOptions := options.Find()

    db, context, cancel := database.GetClient()
    vms := db.Database(database.DefaultDB).Collection(database.VmCollection)
    cur, err := vms.Find(context, bson.D{{"user_id", id}}, findOptions)

    if err != nil {
        defer cancel()
        defer db.Disconnect(context)
        defer cur.Close(context)
        return nil
    }

    for cur.Next(context) {
        var elem map[string]interface{}
        err := cur.Decode(&elem)
        if err != nil {
            fmt.Println("An error occurred while reading data!")
            defer cancel()
            defer db.Disconnect(context)
            defer cur.Close(context)
            return nil
        }

        members := GetVmGroupMembers(elem["server_id"].(string))
        elem["group_members"] = members

        virtualMachines = append(virtualMachines, elem)
    }

    var tmp []database.VirtualMachine

    for _, v := range virtualMachines {
        var t database.VirtualMachine
        d, _ := bson.Marshal(v)
        bson.Unmarshal(d, &t)
        tmp = append(tmp, t)
    }

    var data []database.VirtualMachine

    for _, v := range tmp {
        image := GetImageByImageId(v.ServerImage)
        v.ImageReadRootPassword = image.ImageReadRootPassword
        v.ImageDisplayName = image.ImageDisplayName

        data = append(data, v)
    }

    defer cancel()
    defer db.Disconnect(context)
    defer cur.Close(context)
    return data
}

func GetVMById(id interface{}) (v database.VirtualMachine, erro error) {
    var result database.VirtualMachine
    db, context, cancel := database.GetClient()
    vms := db.Database(database.DefaultDB).Collection(database.VmCollection)
    err := vms.FindOne(context, bson.D{{"server_id", id}}).Decode(&result)

    if err != nil {
        defer cancel()
        defer db.Disconnect(context)
        return result, err
    }

    members := GetVmGroupMembers(result.ServerId)
    result.GroupMembers = members

    image := GetImageByImageId(result.ServerImage)
    result.ImageReadRootPassword = image.ImageReadRootPassword
    result.ImageDisplayName = image.ImageDisplayName

    defer cancel()
    defer db.Disconnect(context)
    return result, nil
}

func UpdateVMStatusById(id interface{}, status string) bool {

    findFilter := bson.M{"server_id": id}
    updateFilter := bson.D{{"$set", bson.M{"server_status": status}}}

    db, context, cancel := database.GetClient()
    vms := db.Database(database.DefaultDB).Collection(database.VmCollection)
    _, err := vms.UpdateMany(context, findFilter, updateFilter)

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

func DeleteVMById(id interface{}) (r int, error error) {

    filter := bson.M{"server_id": bson.M{"$eq": id}}

    db, context, cancel := database.GetClient()
    vms := db.Database(database.DefaultDB).Collection(database.VmCollection)
    res, err := vms.DeleteMany(context, filter)

    if err != nil {
        defer cancel()
        defer db.Disconnect(context)
        return 0, err
    }

    defer cancel()
    defer db.Disconnect(context)
    return int(res.DeletedCount), nil
}

func GetVmGroupMembers(server_id string) []string {
    var virtualMachines []*database.VirtualMachine
    findOptions := options.Find()

    db, context, cancel := database.GetClient()
    vms := db.Database(database.DefaultDB).Collection(database.VmCollection)
    cur, err := vms.Find(context, bson.D{{"server_id", server_id}}, findOptions)

    if err != nil {
        defer cancel()
        defer db.Disconnect(context)
        defer cur.Close(context)
        return nil
    }

    for cur.Next(context) {
        var elem database.VirtualMachine
        err := cur.Decode(&elem)
        if err != nil {
            fmt.Println("An error occurred while reading data!")
            defer cancel()
            defer db.Disconnect(context)
            defer cur.Close(context)
            return nil
        }

        virtualMachines = append(virtualMachines, &elem)
    }

    var members []string
    for _, v := range virtualMachines {
        userName := ""
        if strings.Contains(v.UserId, "@uia.no") {
            userName = strings.Replace(v.UserId, "@uia.no", "", -1)
        } else if strings.Contains(v.UserId, "@student.uia.no") {
            userName = strings.Replace(v.UserId, "@student.uia.no", "", -1)
        }

        members = append(members, userName)
    }

    defer cancel()
    defer db.Disconnect(context)
    defer cur.Close(context)
    return members
}

func CheckIfOwnsVm(server_id interface{}, userId string) bool {
    var result database.VirtualMachine
    db, context, cancel := database.GetClient()
    vms := db.Database(database.DefaultDB).Collection(database.VmCollection)
    err := vms.FindOne(context, bson.D{
        {"server_id", server_id},
        {"user_id", userId},
    }).Decode(&result)

    if err != nil {
        defer cancel()
        defer db.Disconnect(context)
        return false
    }

    defer cancel()
    defer db.Disconnect(context)
    return true
}