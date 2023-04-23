package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"imageCreator/typs"
)

type ImageAddReq struct {
	Url         string `json:"url"`
	ImageDetail string `json:"image_detail"`
}

type ImageGetReq struct {
	Uuid  string `json:"uuid"`
	Page  int    `json:"page"`  // Only In Public Images
	Count int    `json:"count"` //Common
}

type ImageGetResp struct {
	Page   int          `json:"page"` // Only In Public Images
	Images []typs.Image `json:"images"`
}

func ImageAdd(DB *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, ok := c.Get("UUID")
		if !ok {
			c.AbortWithStatus(500)
			return
		}
		var imageAddReq ImageAddReq
		err := c.BindJSON(&imageAddReq)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		image := typs.Image{
			Model: typs.Model{
				ID: uuid.New().String(),
			},
			URL:         imageAddReq.Url,
			ImageDetail: imageAddReq.ImageDetail,
			UserID:      uid.(string),
		}
		err = DB.Create(&image).Error
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		c.String(200, "OK")
	}
}

func ImageGet(DB *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var imageGetReq ImageGetReq
		err := c.BindJSON(&imageGetReq)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		var images []typs.Image
		if imageGetReq.Uuid != "" {
			// 处理指定用户获取图片功能 基于图片所有者UUID 按需加载

			//Preload("User.Image").Preload("Messages.Image").
			err := DB.
				Preload("User").Preload("Messages").
				Where("user_id = ?", imageGetReq.Uuid).
				Order("id ASC").
				Limit(imageGetReq.Count).
				Offset((imageGetReq.Page - 1) * imageGetReq.Count).
				Find(&images).Error
			if err != nil {
				c.AbortWithStatus(500)
				return
			}
			var imageGetResp = ImageGetResp{Images: images, Page: imageGetReq.Page}
			c.JSON(200, imageGetResp)
		} else {
			// 处理用户获取公共图片功能 基于图片主键id 按需加载
			err := DB.
				Preload("User").Preload("Messages").
				Order("id ASC").
				Limit(imageGetReq.Count).
				Offset((imageGetReq.Page - 1) * imageGetReq.Count).
				Find(&images).Error
			if err != nil {
				c.AbortWithStatus(500)
				return
			}
			var imageGetResp = ImageGetResp{Images: images, Page: imageGetReq.Page}
			c.JSON(200, imageGetResp)
		}

	}
}
