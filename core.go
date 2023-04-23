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

type Core struct {
	Engine *gin.Engine
	Config *conf.Config
}

func (c *Core) Run() error {
	err := c.Engine.Run(c.Config.Addr)
	return err
}

func InitializeEngineWithPath(Path string) (*Core, error) {
	config := NewConfWithPath(Path)
	db, err := NewDB(config)
	if err != nil {
		panic(err)
	}
	engine := NewEngineWithConfig()
	handleFunc := NewProxyHandleFunc(config)
	engine = NewRouter(db, config, engine)
	engine = SetEngineWithNoRoute(db, config, engine, handleFunc)
	var core = Core{
		Engine: engine,
		Config: config,
	}
	return &core, nil
}

func InitializeEngineWithReader(reader io.Reader) (*Core, error) {
	config := NewConfWithReader(reader)
	db, err := NewDB(config)
	if err != nil {
		panic(err)
	}
	engine := NewEngineWithConfig()
	handleFunc := NewProxyHandleFunc(config)
	engine = NewRouter(db, config, engine)
	engine = SetEngineWithNoRoute(db, config, engine, handleFunc)
	var core = Core{
		Engine: engine,
		Config: config,
	}
	return &core, nil
}

func NewConfWithPath(Path string) *conf.Config {
	Config, err := conf.ParseConfigFromFile(Path)
	if err != nil {
		panic(err)
	}
	return Config
}

func NewConfWithReader(reader io.Reader) *conf.Config {
	Config, err := conf.ParseConfig(reader)
	if err != nil {
		panic(err)
	}
	return Config
}

func NewDB(Config *conf.Config) (*gorm.DB, error) {
	DB, err := gorm.Open(mysql.Open(Config.Database), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = DB.AutoMigrate(&typs.User{}, &typs.Image{}, &typs.Message{})
	if err != nil {
		return nil, err
	}
	return DB, nil
}

func NewRouter(DB *gorm.DB, Config *conf.Config, r *gin.Engine) *gin.Engine {
	r.POST("/mix", routes.UserMix(DB, Config))
	userRoutes := r.Group("/user")
	{
		userRoutes.Use(middleware.ParseTokenMiddleWare(Config))
		userRoutes.POST("/setdetail", routes.UserSetDetail(DB))
		userRoutes.POST("/getdetail", routes.UserGetDetail(DB))
	}

	imageRoutes := r.Group("/image")
	{
		imageRoutes.Use(middleware.ParseTokenMiddleWare(Config))
		imageRoutes.POST("/add", routes.ImageAdd(DB))
		imageRoutes.POST("/get", routes.ImageGet(DB))
	}

	messageRoutes := r.Group("/message")
	{
		messageRoutes.Use(middleware.ParseTokenMiddleWare(Config))
		messageRoutes.POST("/add", routes.MessageAdd(DB))
		messageRoutes.POST("/get", routes.MessageGet(DB))
	}

	return r
}

func NewEngineWithConfig() *gin.Engine {
	r := gin.Default()
	err := r.SetTrustedProxies([]string{"192.168.88.0/24"})
	if err != nil {
		panic(err)
	}
	return r
}

func NewProxyHandleFunc(Config *conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		request, err := http.NewRequest(c.Request.Method, Config.ServerConfig.Addr+c.Request.URL.Path, c.Request.Body)
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
	}
}

func SetEngineWithNoRoute(DB *gorm.DB, Config *conf.Config, r *gin.Engine, handleFunc gin.HandlerFunc) *gin.Engine {
	r.NoRoute(middleware.ParseTokenMiddleWare(Config),
		middleware.UserBillingBef(DB, Config),
		middleware.UserBillingAft(DB, Config),
		handleFunc,
	)
	return r
}
