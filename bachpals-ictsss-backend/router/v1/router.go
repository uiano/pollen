package v1

import (
    "github.com/gin-gonic/gin"
    "gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/middleware"
)

func Router() *gin.Engine {
    router := gin.Default()

    router.Use(middleware.Cors)
    router.Use(middleware.Json())

    router.GET("/oauth2/provider", HandleSignIn)
    router.GET("/oauth2/redirect", HandleCallback)
    router.POST("/oauth2/userdata", HandleUserdata)
    router.GET("/oauth2/logout", HandleLogout)

    v1 := router.Group("/api/v1")
    {
        admin := v1.Group("/admin")
        {
            admin.GET("/:id", middleware.Authenticate, GetAdministrator)
            admin.GET("/", middleware.Authenticate, ListAdministrators)
            admin.POST("/", middleware.Authenticate, AddAdministrator)
            admin.DELETE("/", middleware.Authenticate, DelAdministrator)
            admin.PUT("/", middleware.Authenticate, UpdateAdministrator)
        }

        images := v1.Group("/image")
        {
            images.GET("/:id", middleware.Authenticate, GetImage)
            images.GET("/", middleware.Authenticate, GetImages)
            images.GET("/published", middleware.Authenticate, GetPublishedImages)
            images.POST("/", middleware.Authenticate, AddImage)
            images.DELETE("/", middleware.Authenticate, DeleteImage)
            images.PUT("/", middleware.Authenticate, UpdateImage)
            images.GET("/server", middleware.Authenticate, GetServerImages)
            images.GET("/config", middleware.Authenticate, GetImagesConfig)
        }

        vms := v1.Group("/vms")
        {
            vms.GET("/", middleware.Authenticate, GetVMs)
            vms.GET("/all", middleware.Authenticate, GetAllVms)
            vms.POST("/", middleware.Authenticate, OrderVM)
            vms.POST("/canvas", middleware.Authenticate, OrderVMFromCanvas)
            vms.POST("/canvas/all", middleware.Authenticate, OrderVMFromCanvasAllStudents)

            // Grouped by VM id
            vms.GET("/:id/status", middleware.Authenticate, StatusVM)
            vms.POST("/:id/start", middleware.Authenticate, StartVM)
            vms.POST("/:id/stop", middleware.Authenticate, StopVM)
            vms.POST("/:id/reboot", middleware.Authenticate, RebootVM)
            vms.POST("/:id/respawn", middleware.Authenticate, RespawnVM)
            vms.DELETE("/:id", middleware.Authenticate, DeleteVM)
            vms.GET("/:id/console", middleware.Authenticate, GenerateConsoleUrl)
            vms.GET("/:id/password", middleware.Authenticate, GetPassword)
        }

        courses := v1.Group("/courses")
        {
            courses.GET("/", middleware.Authenticate, GetCourses)
            courses.GET("/:id/users", middleware.Authenticate, GetCourseStudents)
            courses.GET("/:id/groups", middleware.Authenticate, GetCourseGroups)
            courses.GET("/groups/:id/users", middleware.Authenticate, GetGroupUsers)
        }
    }

    return router
}