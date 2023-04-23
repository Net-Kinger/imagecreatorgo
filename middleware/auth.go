package middleware

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/sha3"
	"imageCreator/conf"
	"strings"
	"time"
)

var headerByte []byte

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	UserID string `json:"user_id"`
	Exp    int64  `json:"exp"`
}

func init() {
	header, err := json.Marshal(&Header{
		Alg: "SHA3-512",
		Typ: "JWT",
	})
	if err != nil {
		panic(err)
	}
	headerByte = header
}

// jwt生成 : Step1: 初始化Alg,Typ部分，填入UserID和Expire过期时间 Step2 将Header,Payload部分分别序列化为Json，将它们使用.进行拼接后转换为特殊的base64编码字符串
// 初始化hmac对象(并指定签名方法和密钥)，并写入hmac对象，h.sum(nil) 获取签名，将Header.Payload.签名base64编码 返回即可

// jwt校验: 使用split对jwtToken用.分割，将[0:2]部分使用.拼接，创建hmac对象(并指定签名方法和密钥)，将[0:2]拼接后写入hmac对象，hmac.Sum计算签名，若签名!=[2]则说明Token被修改过

func GenerateToken(id string, Config *conf.Config) (string, error) {
	payload := Payload{
		UserID: id,
		Exp:    time.Now().Unix() + Config.Auth.ExpireTime,
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	hp := base64.URLEncoding.EncodeToString([]byte(strings.Join([]string{string(headerByte), string(payloadJson)}, ".")))
	hc := hmac.New(sha3.New512, []byte(Config.Auth.Secret))
	n, err := hc.Write([]byte(hp))
	if err != nil || n == 0 {
		return "", nil
	}
	sum := hc.Sum(nil)
	sumBase64 := base64.URLEncoding.EncodeToString(sum)
	return strings.Join([]string{base64.URLEncoding.EncodeToString(headerByte), base64.URLEncoding.EncodeToString(payloadJson), sumBase64}, "."), nil
}

func ParseTokenMiddleWare(Config *conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := c.GetHeader("Authorization")
		hps := strings.Split(jwtToken, ".")
		if len(hps) != 3 {
			c.AbortWithStatus(500)
			return
		}
		header, err := base64.URLEncoding.DecodeString(hps[0])
		payl, err := base64.URLEncoding.DecodeString(hps[1])
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		hpp := strings.Join([]string{string(header), string(payl)}, ".")
		hp := base64.URLEncoding.EncodeToString([]byte(hpp))
		hc := hmac.New(sha3.New512, []byte(Config.Auth.Secret))
		_, err = hc.Write([]byte(hp))
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		sum := hc.Sum(nil)
		sumStr := base64.URLEncoding.EncodeToString(sum)
		if sumStr != hps[2] {
			c.AbortWithStatus(500)
			return
		}
		payloadByte, err := base64.URLEncoding.DecodeString(hps[1])
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		var payload Payload
		err = json.Unmarshal(payloadByte, &payload)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		if time.Now().Unix() >= payload.Exp {
			c.AbortWithStatus(500)
			return
		}
		c.Set("UUID", payload.UserID)
		c.Next()
	}
}
