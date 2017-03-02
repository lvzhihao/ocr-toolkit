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
		if err != nil {
			log.Println(err)
			ctx.JSON(200, map[string]interface{}{"error": err})
		}

		imagick.Initialize()
		defer imagick.Terminate()
		mw := imagick.NewMagickWand()
		err = mw.ReadImage("assert/ocrkit-demo.jpg")
		if err != nil {
			log.Println(err)
			ctx.JSON(200, map[string]interface{}{"error": err})
		}

		//mw.ThresholdImageChannel(imagick.CHANNEL_RED, -0.5)
		//mw.ThresholdImageChannel(imagick.CHANNEL_BLUE, 50.00)
		//mw.ThresholdImageChannel(imagick.CHANNEL_GREEN, 50.00)
		pw := imagick.NewPixelWand()
		pw.SetColor("gray")
		mw.WhiteThresholdImage(pw)
		mw.SetImageClipMask(mw)

		//mw.SetColorspace(imagick.COLORSPACE_GRAY)
		//mw.SetImageClipMask(mw)

		//		rectangleKi := imagick.NewKernelInfoBuiltIn(imagick.KERNEL_RECTANGLE, "3x4")
		//		mw.MorphologyImage(imagick.MORPHOLOGY_CLOSE, 1, rectangleKi)
		//		mw.SetImageClipMask(mw)
		//
		/*
			pw := imagick.NewPixelWand()
			pw.SetAlpha(1.0)
			pw.SetColor("white")
			mw.SetImageBackgroundColor(pw)
			mw.SetColorspace(imagick.COLORSPACE_GRAY)
		*/

		//ki := imagick.NewKernelInfoBuiltIn(imagick.KERNEL_OCTAGON, "3")
		//mw.MorphologyImage(imagick.MORPHOLOGY_SMOOTH, 3, ki)
		//ki := imagick.NewKernelInfoBuiltIn(imagick.KERNEL_SQUARE, "1")
		//mw.MorphologyImage(imagick.MORPHOLOGY_CLOSE, 2, ki)

		squareKi := imagick.NewKernelInfoBuiltIn(imagick.KERNEL_SQUARE, "")
		mw.MorphologyImage(imagick.MORPHOLOGY_ERODE, 1, squareKi)
		mw.SetImageClipMask(mw)

		//mw.ThresholdImage(0)
		//mw.SetImageClipMask(mw)

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
