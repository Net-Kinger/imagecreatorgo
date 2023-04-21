package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"imageCreator/db"
)

type ImageAddReq struct {
	Url         string `json:"url"`
	ImageDetail string `json:"imageDetail"`
}

type ImageGetReq struct {
	Uuid  string `json:"uuid"`
	Page  int    `json:"page"`
	Count int    `json:"count"`
}

type ImageGetResp struct {
	Page   int        `json:"page"`
	Count  int        `json:"count"`
	Images []db.Image `json:"images"`
}

func ImageAdd() func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, ok := c.Get("UUID")
		if !ok {
			c.AbortWithError(500, errors.New("images.30Line Error"))
		}
		var imageAddReq ImageAddReq
		err := c.BindJSON(&imageAddReq)
		if err != nil {
			c.AbortWithError(500, err)
		}
		uidImg := uuid.New().String()
		image := db.Image{
			URL:         imageAddReq.Url,
			ImageDetail: imageAddReq.ImageDetail,
			Uuid:        uidImg,
			User: db.User{
				Uuid: uid.(string),
			},
		}
		err = db.DB.Create(&image).Error
		if err != nil {
			c.AbortWithError(500, err)
		}
		c.String(200, "OK")
	}
}

func ImageGet() func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}
