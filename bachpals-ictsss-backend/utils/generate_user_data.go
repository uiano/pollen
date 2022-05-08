package utils

import (
    "bytes"
    "fmt"
    "github.com/spf13/viper"
    "io/ioutil"
    "strings"
)

func readFile(templatePath string) []byte {
    input, err := ioutil.ReadFile(templatePath)
    if err != nil {
        fmt.Println(err)
    }

    return input
}

// users has to be a comma separated string of usernames or just a single user without a comma at the end.
func generateSudoers(users string) []byte {
    spittedUsers := strings.Split(users, ",")

    var tmp []byte

    for k, v := range spittedUsers {
        t := ""
        if k == 0 {
            t = fmt.Sprintf("%s    ALL=(ALL:ALL) ALL\n", v)
        } else {
            t = fmt.Sprintf("    %s    ALL=(ALL:ALL) ALL\n", v)
        }
        tmp = append(tmp, []byte(t)...)
    }

    return tmp
}

// From a comma separated list of users, removes all occurrences of @uia.no or @student.uia.no
func CleanUserNames(users string) string {
    spittedUsers := strings.Split(users, ",")

    var tmp []string

    for _, v := range spittedUsers {
        if strings.Contains(v, "@uia.no") {
            tmp = append(tmp, strings.Replace(v, "@uia.no", "", -1))
        } else if strings.Contains(v, "@student.uia.no") {
            tmp = append(tmp, strings.Replace(v, "@student.uia.no", "", -1))
        }
    }
    return strings.Join(tmp, ",")
}

// GenerateUserData
// image is name of the image, deb, kali, ikt207 etc.
// users has to be a comma separated string of usernames or just a single user without a comma at the end.
func GenerateUserData(imageConfig string, users string) []byte {
    if len(imageConfig) == 0 {
        fmt.Println("Image needs to have a valid value!")
        return nil
    }

    if len(users) == 0 {
        fmt.Println("Users needs to have a valid value!")
        return nil
    }

    cleanUsers := CleanUserNames(users)

    sssd := readFile(viper.GetString("IKT_STACK_TEMPLATES_CONFIGS_DIR") + viper.GetString("IKT_STACK_SSSD_TEMPLATE_NAME"))
    sssdOutput := bytes.Replace(sssd, []byte("{USERS}"), []byte(cleanUsers), -1)

    userData := readFile(viper.GetString("IKT_STACK_TEMPLATES_USERDATA_DIR") + imageConfig)

    t := ""

    tmp := strings.Split(string(sssdOutput), "\n")
    for k, v := range tmp {
        if k == 0 {
            t += fmt.Sprintf("%s\n", v)
        } else {
            t += fmt.Sprintf("  %s\n", v)
        }
    }

    userDataOutput := bytes.Replace(userData, []byte("{SSSD_CONF}"), []byte(t), -1)
    userDataOutput = bytes.Replace(userDataOutput, []byte("{SUDOERS}"), generateSudoers(cleanUsers), -1)

    return userDataOutput
}