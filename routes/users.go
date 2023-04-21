package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"imageCreator/auth"
	"imageCreator/conf"
	"imageCreator/db"
	"net/http"
)

type UserMixReq struct {
	IsRegistry  bool   `json:"isRegistry"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserMixResp struct {
	Uuid       string `json:"uuid"`
	JwtToken   string `json:"jwt_token"`
	Tokens     int64  `json:"tokens"`
	ExpireTime int64  `json:"expire_time"`
}

type UserDetail struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	NewPassword string `json:"new_password"`
	Token       int64  `json:"token"`
}

func UserMix() func(c *gin.Context) {
	return func(c *gin.Context) {
		var userMixReq UserMixReq
		err := c.Bind(&userMixReq)
		if err != nil {
			c.String(200, "Error")
		}

		if !userMixReq.IsRegistry {
			// 处理登陆逻辑
			var user db.User
			err := db.DB.Find(&user).Where("phone_number = ?", userMixReq.PhoneNumber).Error
			if err != nil {
				c.String(500, "users 49 Line:"+err.Error())
				return
			}
			if user.Password != userMixReq.Password {
				c.String(500, "密码错误")
				return
			}
			jwtToken, err := auth.GenerateToken(user.Uuid)
			if err != nil {
				c.String(500, "jwt生成错误"+err.Error())
				return
			}
			userMixResp := UserMixResp{
				Uuid:       user.Uuid,
				JwtToken:   jwtToken,
				Tokens:     user.Token,
				ExpireTime: conf.Conf.Auth.ExpireTime,
			}
			c.JSON(http.StatusOK, userMixResp)
			return
		}
		// 处理注册逻辑
		uid := uuid.New().String()
		err = db.DB.Create(&db.User{
			PhoneNumber: userMixReq.PhoneNumber,
			Password:    userMixReq.Password,
			Uuid:        uid,
			Token:       0,
			Images:      nil,
			Messages:    nil,
		}).Where("phone_number <> ?", userMixReq.PhoneNumber).Error
		if err != nil {
			c.String(200, "users 55 Line:"+err.Error())
			return
		}

		jwtToken, err := auth.GenerateToken(uid)
		if err != nil {
			c.String(200, "users 61 Line:"+err.Error())
			return
		}
		userMixResp := UserMixResp{
			Uuid:       uid,
			JwtToken:   jwtToken,
			Tokens:     0,
			ExpireTime: conf.Conf.Auth.ExpireTime,
		}
		c.JSON(http.StatusOK, userMixResp)
	}
}

func UserSetDetail() func(c *gin.Context) {
	return func(c *gin.Context) {
		var userDetail UserDetail
		if err := c.BindJSON(&userDetail); err != nil {
			c.String(500, "users.98:"+err.Error())
			return
		}
		uid, ok := c.Get("UUID")
		if !ok {
			c.String(500, "users.102")
			return
		}
		var user db.User
		err := db.DB.Where("uuid = ?", uid).Find(&user).Error
		if err != nil {
			c.String(500, "用户不存在")
			return
		}

		user.PhoneNumber = userDetail.PhoneNumber
		user.Password = userDetail.NewPassword
		user.Name = userDetail.Name
		err = db.DB.Where("uuid = ?", uid).Updates(&user).Error
		if err != nil {
			c.String(500, "user.111:"+err.Error())
			return
		}
		c.String(200, "OK")
	}
}

func UserGetDetail() func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, ok := c.Get("UUID")
		if !ok {
			c.String(500, "users.121")
			return
		}
		var user db.User
		err := db.DB.Where("uuid = ?", uid).Find(&user).Error
		if err != nil {
			c.String(500, "用户不存在"+err.Error())
			return
		}

		var userDetail = UserDetail{
			Token:       user.Token,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		}

		c.JSON(200, &userDetail)
	}
}
