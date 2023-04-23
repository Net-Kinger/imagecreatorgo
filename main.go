package main

import "github.com/gin-gonic/gin"

func main() {
	gin.SetMode(gin.ReleaseMode)
	Path := "conf/config.yaml"
	engine, err := InitializeEngineWithPath(Path)
	if err != nil {
		panic(err)
	}
	err = engine.Run()
	if err != nil {
		panic(err)
	}
}

//func resolvePNGToBase64Jpeg(img []byte) ([]byte, error) {
//	//file, _ := os.Open("1.png")
//	//buf := new(bytes.Buffer)
//	//buf.ReadFrom(file)
//	//jpeg, _ := resolvePNGToBase64Jpeg(buf.Bytes())
//	//fileNew, _ := os.Create("2.jpeg")
//	//fileNew.Write(jpeg)
//
//	buf := bytes.NewBuffer(img)
//	//decodeString, _ := base64.StdEncoding.DecodeString(img)
//	image, _ := png.Decode(buf)
//	out := new(bytes.Buffer)
//	jpeg.Encode(out, image, &jpeg.Options{Quality: 75})
//	return out.Bytes(), nil
//}
