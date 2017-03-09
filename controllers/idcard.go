package controllers

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/otiai10/gosseract"
	"gopkg.in/gographics/imagick.v2/imagick"
	"gopkg.in/kataras/iris.v6"
)

func IDCardApiOptions(ctx *iris.Context) {
	ctx.SetHeader("Access-Control-Allow-Origin", "*")
	ctx.HTML(200, "")
}

func IDCardApi(ctx *iris.Context) {

	var image map[string]string
	err := ctx.ReadJSON(&image)
	if err != nil {
		log.Println(err)
		ctx.JSON(200, map[string]string{"error": "image error"})
	} else {
		var imageData []byte
		var err error
		imageData, err = base64.StdEncoding.DecodeString(image["image"])
		var whitelist = "0123456789xX"
		if err != nil {
			log.Println(err)
			ctx.JSON(200, map[string]interface{}{"error": err})
		}
		Lk.Lock()
		defer Lk.Unlock()
		if err != nil {
			log.Println(err)
			ctx.JSON(200, map[string]interface{}{"error": err})
		}

		imagick.Initialize()
		defer imagick.Terminate()
		mw := imagick.NewMagickWand()
		err = mw.ReadImageBlob(imageData)
		if err != nil {
			log.Println(err)
			ctx.JSON(200, map[string]interface{}{"error": err})
		}
		defer mw.Destroy()

		pw := imagick.NewPixelWand()
		pw.SetColor("gray")
		mw.WhiteThresholdImage(pw)
		mw.SetImageClipMask(mw)

		mw.SetType(imagick.IMAGE_TYPE_GRAYSCALE)
		mw.SetImageColorspace(imagick.COLORSPACE_GRAY)
		mw.SetImageClipMask(mw)

		rectangleKi := imagick.NewKernelInfoBuiltIn(imagick.KERNEL_RECTANGLE, "3x1:1,0,1")
		mw.MorphologyImage(imagick.MORPHOLOGY_CLOSE, 2, rectangleKi)
		mw.SetImageClipMask(mw)

		mw.SharpenImage(4.0, 1.5)
		mw.SigmoidalContrastImage(true, 1.8, 10.0)
		mw.SetImageClipMask(mw)

		err = mw.WriteImage("assert/ocrkit-demo.jpg")

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
