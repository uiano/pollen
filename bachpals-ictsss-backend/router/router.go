package router

import (
    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
    v1 "gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/router/v1"
)

func Router() *gin.Engine {
    var r *gin.Engine

    if viper.GetString("IKT_STACK_API_VERSION") == "v1" {
        r = v1.Router()
    }

    return r
}