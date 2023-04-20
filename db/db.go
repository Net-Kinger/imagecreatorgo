package db

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"imageCreator/conf"
	"math"
	"net/http"
	"strconv"
	"time"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	DB, err = gorm.Open(mysql.Open(conf.Conf.Database), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&User{}, &Image{}, &Message{})
	if err != nil {
		panic(err)
	}
}

func UserBillingBef() func(c *gin.Context) {
	return func(c *gin.Context) {
		uuid, ok := c.Get("UUID")
		if ok == false {
			c.String(http.StatusOK, "鉴权错误")
		}
		var user = User{}
		err := DB.Find(&user, "uuid = ?", uuid).Error
		if err != nil {
			c.String(http.StatusOK, "UUID未找到")
		}
		if user.Token <= conf.Conf.TokenRelation.MinToken {
			c.String(http.StatusOK, "账户Token低于"+strconv.FormatInt(conf.Conf.TokenRelation.MinToken, 10)+"不足,请充值")
		}
		c.Set("User", user)
		c.Next()
	}
}

func UserBillingAft() func(c *gin.Context) {
	return func(c *gin.Context) {
		now := time.Now()
		c.Next()
		i, ok := c.Get("User")
		if !ok {
			c.String(http.StatusOK, "上下文错误")
		}
		user, ok := i.(User)
		if !ok {
			c.String(http.StatusOK, "断言错误")
		}
		consumeCountFloat := time.Since(now).Seconds() * conf.Conf.TokenRelation.Magnification
		var consumeCount = int64(math.Round(consumeCountFloat))
		user.Token = user.Token - consumeCount
		tx := DB.Model(&User{}).Where("uuid = ?", user.Uuid).Updates(&user)
		if tx.Error != nil {
			c.String(http.StatusOK, "数据库异常")
		}
	}
}

func Save() {
	//DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	panic(err)
	//}
	//
	//user := User{}
	//DB.Preload("Images.Messages.User").Preload("Messages.User").Find(&user, 1)
	//fmt.Println(user.Messages[0].User.Name)
	////
	////if err != nil {
	////	panic(err)
	////}
	////
	////user1 := User{
	////	Name: "丽华",
	////	Uuid: "123456",
	////}
	////err = DB.Create(&user1).Error
	////if err != nil {
	////	panic(err)
	////}
	////
	////image1 := Image{
	////	URL:  "https://1.1.1.1:1111/1.jpg",
	////	User: user1,
	////}
	////
	////err = DB.Create(&image1).Error
	////if err != nil {
	////	panic(err)
	////}
	////
	////err = DB.Create(&Message{
	////	Text:  "你好",
	////	User:  user1,
	////	Image: image1,
	////}).Error
	////if err != nil {
	////	panic(err)
	////}
}
