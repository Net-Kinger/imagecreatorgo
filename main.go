package main

import (
	"github.com/gin-gonic/gin"
	"imageCreator/auth"
	"imageCreator/conf"
	"imageCreator/db"
	"time"
)

func main() {
	err := conf.ParseConfigFromFile("conf/config.yaml")
	if err != nil {
		panic(err)
	}

	db.InitDatabase()
	engine := gin.Default()

	engine.Use(auth.ParseTokenMiddleWare())
	engine.GET("/test", db.UserBillingBef(), db.UserBillingAft(), func(context *gin.Context) {
		time.Sleep(10 * time.Second)
	})
	err = engine.Run(conf.Conf.Addr)
	if err != nil {
		return
	}
}
