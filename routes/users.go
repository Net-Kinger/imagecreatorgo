package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"imageCreator/conf"
	"imageCreator/middleware"
	"imageCreator/typs"
	"net/http"
)

type UserMixReq struct {
	IsRegistry  bool   `json:"is_registry"`
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
	Uuid        string `json:"uuid"` // GET
	PhoneNumber string `json:"phone_number"`
	NewPassword string `json:"new_password"` // SET
	Token       int64  `json:"token"`        // GET
}

func UserMix(DB *gorm.DB, Config *conf.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		var userMixReq UserMixReq
		err := c.Bind(&userMixReq)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		if !userMixReq.IsRegistry {
			// 处理登陆逻辑
			var user typs.User
			err := DB.Where("phone_number = ?", userMixReq.PhoneNumber).Find(&user).Error
			if err != nil {
				c.AbortWithStatus(500)
				return
			}
			if user.Password != userMixReq.Password {
				c.AbortWithStatus(500)
				return
			}
			jwtToken, err := middleware.GenerateToken(user.ID, Config)
			if err != nil {
				c.AbortWithStatus(500)
				return
			}
			userMixResp := UserMixResp{
				Uuid:       user.ID,
				JwtToken:   jwtToken,
				Tokens:     user.Token,
				ExpireTime: Config.Auth.ExpireTime,
			}
			c.JSON(http.StatusOK, userMixResp)
			return
		}
		// 处理注册逻辑
		uid := uuid.New()
		err = DB.Create(&typs.User{
			PhoneNumber: userMixReq.PhoneNumber,
			Password:    userMixReq.Password,
			Model: typs.Model{
				ID: uid.String(),
			},
			Token:    Config.TokenRelation.MinToken,
			Images:   nil,
			Messages: nil,
		}).Where("phone_number <> ?", userMixReq.PhoneNumber).Error
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		jwtToken, err := middleware.GenerateToken(uid.String(), Config)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		userMixResp := UserMixResp{
			Uuid:       uid.String(),
			JwtToken:   jwtToken,
			Tokens:     Config.TokenRelation.MinToken,
			ExpireTime: Config.Auth.ExpireTime,
		}
		c.JSON(http.StatusOK, userMixResp)
	}
}

func UserSetDetail(DB *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var userDetail UserDetail
		if err := c.BindJSON(&userDetail); err != nil {
			c.AbortWithStatus(500)
			return
		}
		uid, ok := c.Get("UUID")
		if !ok {
			c.AbortWithStatus(500)
			return
		}

		err := DB.Model(&typs.User{}).Where("id = ?", uid).Updates(map[string]interface{}{
			"PhoneNumber": userDetail.PhoneNumber,
			"Password":    userDetail.NewPassword,
			"Name":        userDetail.Name,
		}).Error
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		c.String(200, "OK")
	}
}

func UserGetDetail(DB *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, ok := c.Get("UUID")
		if !ok {
			c.AbortWithStatus(500)
			return
		}
		var user typs.User
		err := DB.Where("id = ?", uid).Find(&user).Error
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		var userDetail = UserDetail{
			Token:       user.Token,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Uuid:        user.ID,
		}

		c.JSON(200, &userDetail)
	}
}
