package main

import (
	"imageCreator/auth"
	"imageCreator/conf"
	"os"
)

func main() {
	err := conf.ParseConfigFromFile("conf/config.yaml")
	if err != nil {
		panic(err)
	}
	//db.InitDatabase()

	token, _ := auth.GenerateToken("u1")
	create, _ := os.Create("token.txt")
	create.Write([]byte(token))
	create.Close()
	//db.InitDatabase()
	//engine := gin.Default()
	//
	//engine.Use(auth.ParseTokenMiddleWare())
	//engine.GET("/test", db.UserBillingBef(), func(context *gin.Context) {
	//
	//}, db.UserBillingAft())
	//engine.Run(conf.Conf.Addr)
}
