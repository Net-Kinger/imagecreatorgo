package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"imageCreator/conf"
	"imageCreator/middleware"
	"imageCreator/routes"
	"imageCreator/typs"
	"io"
	"net/http"
)

func InitialConf(Path string) *conf.Config {
	Config, err := conf.ParseConfigFromFile(Path)
	if err != nil {
		panic(err)
	}
	return Config
}

func InitialDB(Config *conf.Config) *gorm.DB {
	DB, err := gorm.Open(mysql.Open(conf.Conf.Database), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&typs.User{}, &typs.Image{}, &typs.Message{})
	if err != nil {
		panic(err)
	}
	return DB
}

func InitialRoutes(DB *gorm.DB, Config *conf.Config) gin.RouterGroup {
	var r = gin.RouterGroup{}
	r.POST("/mix", routes.UserMix(DB, Config))
	userRoutes := r.Group("/user")
	{
		userRoutes.Use(middleware.ParseTokenMiddleWare())
		userRoutes.POST("/setdetail", routes.UserSetDetail(DB))
		userRoutes.POST("/getdetail", routes.UserGetDetail(DB))
	}

	imageRoutes := r.Group("/image")
	{
		imageRoutes.Use(middleware.ParseTokenMiddleWare())
		imageRoutes.POST("/add", routes.ImageAdd(DB))
		imageRoutes.POST("/get", routes.ImageGet(DB))
	}

	messageRoutes := r.Group("/message")
	{
		messageRoutes.Use(middleware.ParseTokenMiddleWare())
		messageRoutes.POST("/add", routes.MessageAdd(DB))
		messageRoutes.POST("/get", routes.MessageGet(DB))
	}
	return r
}

func InitialEngine(DB *gorm.DB, Config *conf.Config) *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	err := r.SetTrustedProxies([]string{"192.168.88.0/24"})
	if err != nil {
		panic(err)
	}
	r.NoRoute(middleware.ParseTokenMiddleWare(),
		middleware.UserBillingBef(DB, Config),
		middleware.UserBillingAft(DB, Config),
		func(c *gin.Context) {
			request, err := http.NewRequest(c.Request.Method, conf.Conf.ServerConfig.Addr+c.Request.URL.Path, c.Request.Body)
			if err != nil {
				c.AbortWithStatus(500)
				return
			}
			request.Header = c.Request.Header
			client := &http.Client{}
			resp, err := client.Do(request)
			if err != nil {
				c.AbortWithStatus(500)
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					c.AbortWithStatus(500)
					return
				}
			}(resp.Body)
			for k, v := range resp.Header {
				c.Writer.Header().Set(k, v[0])
			}
			c.Status(200)
			_, err = io.Copy(c.Writer, resp.Body)
			if err != nil {
				c.AbortWithStatus(500)
				return
			}
		})

	return r
}
