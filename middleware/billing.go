package middleware

import (
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"imageCreator/conf"
	"imageCreator/typs"
	"io"
	"math"
	"regexp"
	"time"
)

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

type CoreResponseWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (c CoreResponseWriter) Write(b []byte) (int, error) {
	c.buf.Write(b)
	return c.ResponseWriter.Write(b)
}

func UserBillingAft(DB *gorm.DB, Config *conf.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		writer := CoreResponseWriter{
			ResponseWriter: c.Writer,
			buf:            bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer

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

		t, err := parseTime(writer.buf.Bytes(), true)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		since := time.Since(t)
		ceil := math.Ceil(since.Seconds())
		tx := DB.Model(&user).Update("Token", gorm.Expr("Token - ?", ceil))
		if tx.Error != nil {
			c.AbortWithStatus(500)
			return
		}
	}
}

func parseTime(input []byte, needCompress bool) (time.Time, error) {
	// gzip 解压缩 -> 正则表达式获取timeStamp -> ParseTime
	buf := bytes.NewBuffer(input)
	var byteS []byte
	var err error
	if needCompress {
		reader, err := gzip.NewReader(buf)
		if err != nil {
			return time.Time{}, err
		}
		byteS, err = io.ReadAll(reader)
		if err != nil {
			return time.Time{}, err
		}
	} else {
		byteS, err = io.ReadAll(buf)
		if err != nil {
			return time.Time{}, err
		}
	}
	cp, err := regexp.Compile(`"job_timestamp\\": \\"(\d+)\\"`)
	if err != nil {
		return time.Time{}, err
	}
	subMatch := cp.FindStringSubmatch(string(byteS))
	if len(subMatch) != 2 {
		return time.Time{}, errors.New("正则匹配失败")
	}
	jobStartTime, err := time.ParseInLocation("20060102150405", subMatch[1], time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return jobStartTime, nil
}
