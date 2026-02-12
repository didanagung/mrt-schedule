package main

import (
	"github.com/didanagung/mrt-schedule/config"
	"github.com/didanagung/mrt-schedule/modules/station"
	"github.com/gin-gonic/gin"
)

func main() {
	InitiateRouter()
}

func InitiateRouter() {
	var (
		router = gin.Default()
		api    = router.Group("/v1/api")
	)

	port := config.GetString("PORT")

	// bungkus ke dalam api
	station.Initiate(api)

	router.Run(":" + port)
}
