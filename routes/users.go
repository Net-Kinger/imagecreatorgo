package routes

import (
	"github.com/gin-gonic/gin"
)

type UserMixReq struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserMixResp struct {
	JwtToken   string `json:"jwt_token"`
	Tokens     int    `json:"tokens"`
	ExpireTime string `json:"expire_time"`
}

type UserDetailResp struct {
	PhoneNumber string
	Token       int64
}

type UserSetPasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func UserMix() func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}

func UserDetail() func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}

func UserSetPassword() func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}
