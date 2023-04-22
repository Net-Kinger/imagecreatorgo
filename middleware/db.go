package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"imageCreator/conf"
	"imageCreator/typs"
	"math"
	"time"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	DB, err = gorm.Open(mysql.Open(conf.Conf.Database), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&typs.User{}, &typs.Image{}, &typs.Message{})
	if err != nil {
		panic(err)
	}
}

type Middle struct {
	conf *conf.Config
	db   *gorm.DB
}

func UserBillingBef(DB *gorm.DB, Config *conf.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		uuid, ok := c.Get("UUID")
		if ok == false {
			c.AbortWithStatus(500)
			return
		}
		var user = typs.User{}
		err := DB.Find(&user, "id = ?", uuid).Error
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		if user.Token < Config.TokenRelation.MinToken {
			c.AbortWithStatus(500)
			return
		}
		c.Set("User", user)
		c.Next()
	}
}

func UserBillingAft(DB *gorm.DB, Config *conf.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		now := time.Now()
		c.Next()
		i, ok := c.Get("User")
		if !ok {
			c.AbortWithStatus(500)
			return
		}
		user, ok := i.(typs.User)
		if !ok {
			c.AbortWithStatus(500)
			return
		}
		consumeCountFloat := time.Since(now).Seconds() * Config.TokenRelation.Magnification
		var consumeCount = int64(math.Round(consumeCountFloat))
		user.Token = user.Token - consumeCount
		tx := DB.Model(&typs.User{}).Where("id = ?", user.ID).Updates(&user)
		if tx.Error != nil {
			c.AbortWithStatus(500)
			return
		}
	}
}
