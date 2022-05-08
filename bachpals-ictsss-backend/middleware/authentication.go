package middleware

import (
    "encoding/json"
    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
    "gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/httputils"
    "io"
    "net/http"
    "strings"
)

func Authenticate(c *gin.Context) {
    authorization := c.GetHeader("Authorization")

    if len(authorization) <= 0 {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    client := &http.Client{}

    req, err := http.NewRequest("GET", viper.GetString("IKT_STACK_AUTH_MIDDLEWARE_URL"), nil)

    if err != nil {
        httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while reaching endpoint", err)
        return
    }

    req.Header.Add("Authorization", authorization)
    resp, err := client.Do(req)

    // Any status code that is not 200 OK should return false
    if resp.StatusCode == http.StatusUnauthorized {
        httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Missing authorization header", err)
        return
    }

    if resp.StatusCode != 200 {
        httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Something went wrong!", err)
        return
    }

    // Read users email address from Feide, then pass it to next request.
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while reading data", err)
        return
    }

    var jsonData map[string]interface{}
    err = json.Unmarshal(bodyBytes, &jsonData)
    if err != nil {
        httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while encoding data", err)
        return
    }

    user := strings.Split(jsonData["email"].(string), "@")
    formattedUser := user[0] + "@uia.no"

    c.Set("user_id", formattedUser)

    c.Next()
}