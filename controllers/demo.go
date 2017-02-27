package controllers

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/otiai10/gosseract"
	"gopkg.in/gographics/imagick.v2/imagick"
	"gopkg.in/kataras/iris.v6"
)

var Client *gosseract.Client
var Lk sync.Mutex

func init() {
	// Using client
	var err error
	Client, err = gosseract.NewClient()
	if err != nil {
		log.Fatal(err)
	}
}

func Demo(ctx *iris.Context) {
	ctx.Render("demo/index.html", nil)
}

func DemoApi(ctx *iris.Context) {

	var image map[string]string
	err := ctx.ReadJSON(&image)
	if err != nil {
		log.Println(err)
		ctx.JSON(200, map[string]string{"error": "image error"})
	} else {
		var imageData []byte
		var err error
		imageData, err = base64.StdEncoding.DecodeString(image["image"])
		var whitelist = ""
		if _, ok := image["whitelist"]; ok {
			whitelist = image["whitelist"]
		}
		log.Println(whitelist)
		if err != nil {
			log.Println(err)
			ctx.JSON(200, map[string]interface{}{"error": err})
		}
		Lk.Lock()
		defer Lk.Unlock()
		err = ioutil.WriteFile("assert/ocrkit-demo.jpg", imageData, 0600)

		imagick.Initialize()
		defer imagick.Terminate()
		mw := imagick.NewMagickWand()
		err = mw.ReadImage("assert/ocrkit-demo.jpg")
		if err != nil {
			log.Println(err)
			ctx.JSON(200, map[string]interface{}{"error": err})
		}
		mw.SharpenImage(4.0, 1.5)
		mw.SigmoidalContrastImage(true, 1.8, 10.0)
		err = mw.WriteImage("assert/ocrkit-demo.jpg")
		if err != nil {
			log.Println(err)
			ctx.JSON(200, map[string]interface{}{"error": err})
		}

		if err != nil {
			log.Println(err)
			ctx.JSON(200, map[string]interface{}{"error": err})
		}
		var out string
		if whitelist == "" {
			out = gosseract.Must(gosseract.Params{
				Src:       "assert/ocrkit-demo.jpg",
				Languages: "eng+chi_sim",
			})
		} else {
			out = gosseract.Must(gosseract.Params{
				Src:       "assert/ocrkit-demo.jpg",
				Whitelist: whitelist,
				Languages: "eng+chi_sim",
			})
		}
		log.Println(out)
		ctx.JSON(200, map[string]interface{}{
			"data": strings.Replace(strings.TrimSpace(out), " ", "", -1),
		})
	}
}
