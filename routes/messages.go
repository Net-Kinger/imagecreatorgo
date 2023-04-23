package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"imageCreator/typs"
	"time"
)

type MessageAddReq struct {
	Text      string `json:"text"`
	ImageUUID string `json:"image_uuid"`
	DstUUID   string `json:"dst_uuid"`
}

type MessageGetReq struct {
	UUID string `json:"image_uuid"` // 图片UUID
}

type MessageGetResp struct {
	UserID   string
	DstID    string
	Text     string
	UserName string
	DstName  string
	Time     time.Time
}

func MessageAdd(DB *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, ok := c.Get("UUID")
		if !ok {
			c.AbortWithStatus(500)
			return
		}
		var messageAddReq MessageAddReq
		err := c.BindJSON(&messageAddReq)
		if err != nil {
			c.Abort()
			c.String(500, "JSON错误")
			return
		}
		var message = typs.Message{
			Model:   typs.Model{ID: uuid.New().String()},
			Text:    messageAddReq.Text,
			UserID:  uid.(string),
			DstID:   messageAddReq.DstUUID,
			ImageID: messageAddReq.ImageUUID,
		}
		err = DB.Create(&message).Error
		if err != nil {
			c.Abort()
			c.String(500, err.Error())
			return
		}
		c.String(200, "OK")
	}
}

func MessageGet(DB *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var messageGetReq MessageGetReq
		err := c.BindJSON(&messageGetReq)
		if err != nil {
			c.Abort()
			c.String(500, err.Error())
			return
		}
		var message []typs.Message
		err = DB.Preload("User").Preload("Dst").Where("image_id = ?", messageGetReq.UUID).Find(&message).Error
		if err != nil {
			c.Abort()
			c.String(500, err.Error())
			return
		}
		var messageRespList = make([]MessageGetResp, 0, 10)
		for i := 0; i < len(message); i++ {
			messageResp := MessageGetResp{
				UserID:   message[i].UserID,
				DstID:    message[i].DstID,
				UserName: message[i].User.Name,
				DstName:  message[i].Dst.Name,
				Time:     message[i].CreatedAt,
				Text:     message[i].Text,
			}
			messageRespList = append(messageRespList, messageResp)
		}
		c.JSON(200, &messageRespList)
	}
}
