package database

import (
    "time"
)

const VirtualMachineStatusInactive = "SHUTOFF"
const VirtualMachineStatusActive = "ACTIVE"

const ServerStatusPollingTime = 300 // Seconds

type VirtualMachine struct {
    ServerIp     string    `bson:"server_ip"`
    ServerImage  string    `bson:"server_image"`
    ServerName   string    `bson:"server_name"`
    ServerStatus string    `bson:"server_status"`
    UserId       string    `bson:"user_id"`
    ServerId     string    `bson:"server_id"`
    Created      time.Time `bson:"created"`
    GroupMembers []string  `bson:"group_members"`
    VirtualMachineImageMeta
}

type VirtualMachineImageMeta struct {
    ImageReadRootPassword bool   `bson:"image_read_root_password"`
    ImageDisplayName      string `bson:"image_display_name"`
}

type Admin struct {
    UserId string `bson:"user_id"`
    Name   string `bson:"name"`
}

type Application struct {
    FirstRun int `bson:"first_run"`
}

type Images struct {
    Id                    string `bson:"_id"`
    Published             string `bson:"published"`
    ImageId               string `bson:"image_id"`
    ImageName             string `bson:"image_name"`
    ImageDescription      string `bson:"image_description"`
    ImageDisplayName      string `bson:"image_display_name"`
    ImageConfig           string `bson:"image_config"`
    ImageReadRootPassword bool   `bson:"image_read_root_password"`
}