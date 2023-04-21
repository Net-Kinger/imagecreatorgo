package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"image/jpeg"
	"image/png"
	"imageCreator/auth"
	"imageCreator/conf"
	"imageCreator/db"
	"imageCreator/routes"
)

func main() {
	err := conf.ParseConfigFromFile("conf/config.yaml")
	if err != nil {
		panic(err)
	}

	db.InitDatabase()
	engine := gin.Default()

	//engine.Use()
	//engine.GET("/test", db.UserBillingBef(), db.UserBillingAft(), func(context *gin.Context) {
	//	time.Sleep(10 * time.Second)
	//})

	engine.POST("/user/userMix", routes.UserMix())
	engine.POST("/user/getDetail", auth.ParseTokenMiddleWare(), routes.UserGetDetail())
	engine.POST("/user/setDetail", auth.ParseTokenMiddleWare(), routes.UserSetDetail())
	err = engine.Run(conf.Conf.Addr)
	if err != nil {
		return
	}

	//file, _ := os.Open("1.png")
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(file)
	//jpeg, _ := resolvePNGToBase64Jpeg(buf.Bytes())
	//fileNew, _ := os.Create("2.jpeg")
	//fileNew.Write(jpeg)
}

func resolvePNGToBase64Jpeg(img []byte) ([]byte, error) {
	buf := bytes.NewBuffer(img)
	//decodeString, _ := base64.StdEncoding.DecodeString(img)
	image, _ := png.Decode(buf)
	out := new(bytes.Buffer)
	jpeg.Encode(out, image, &jpeg.Options{Quality: 75})
	return out.Bytes(), nil
}
