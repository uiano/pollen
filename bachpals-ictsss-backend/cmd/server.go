package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/database/repositories"
	_ "gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/docs"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/router"
)

// @title           ICT-Stack Self-Service API
// @version         1.0
// @description     This is the backend to the ICTSSS service.
// @termsOfService  http://swagger.io/terms/

// @contact.name ICTSSS Support
// @contact.email <TBD>

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3001
// @BasePath  /api/v1
func serve() {
	gin.SetMode(viper.GetString("IKT_STACK_SERVER_MODE"))

	r := router.Router()

	// Swaggo MUST only run when Mode == "debug"
	if viper.GetString("IKT_STACK_SERVER_MODE") == gin.DebugMode {
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	srv := &http.Server{
		Addr:    viper.GetString("IKT_STACK_SERVER_PORT"),
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func main() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		_ = fmt.Errorf("Fatal error config file: %w \n", err)
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)

		val, ok := os.LookupEnv(pair[0])
		if !ok {
			fmt.Printf("%s not set\n", pair[0])
		} else {
			viper.Set(pair[0], val)
		}
	}

	if len(os.Args) > 1 {
		arg := os.Args[1]

		if arg == "--reset" {
			repositories.RemoveInitializationStatus()
			os.Exit(0)
		}

		if arg == "--help" {
			fmt.Println(`Available commands
--reset    Removes all administrators and resets default settings in database.
--init     Initializes administrators and default settings in database.
--serve    Starts the http server.
--help     Shows this page.`)
			os.Exit(0)
		}

		if arg == "--serve" {
			serve()
		}

		if arg == "--init" {
			// If application is not initialized, add default administrator.
			if repositories.CheckInitializationStatus() {
				repositories.InitializeDefaultAdministrator()
			}
		}
	}
}
