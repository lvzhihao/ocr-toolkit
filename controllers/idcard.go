package controllers

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"strings"

	leptonica "github.com/GeertJohan/go.leptonica"
	tesseract "github.com/GeertJohan/go.tesseract"
	"gopkg.in/gographics/imagick.v2/imagick"
	"gopkg.in/kataras/iris.v6"
)

var IDCardTess *tesseract.Tess

func init() {
	tessdata_prefix := os.Getenv("TESSDATA_PREFIX")
	if tessdata_prefix == "" {
		tessdata_prefix = "/usr/local/share"
	}
	IDCardTess, _ = tesseract.NewTess(filepath.Join(tessdata_prefix, "tessdata"), "eng")
}

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

		rectangleKi := imagick.NewKernelInfoBuiltIn(imagick.KERNEL_RECTANGLE, "3x2:1,0,1")
		mw.MorphologyImage(imagick.MORPHOLOGY_CLOSE, 2, rectangleKi)
		mw.SetImageClipMask(mw)

		mw.SharpenImage(4.0, 1.5)
		mw.SigmoidalContrastImage(true, 1.8, 10.0)
		mw.SetImageClipMask(mw)

		if err != nil {
			log.Println(err)
			ctx.JSON(200, map[string]interface{}{"error": err})
		}

		IDCardTess.SetPageSegMode(tesseract.PSM_CIRCLE_WORD)
		IDCardTess.SetVariable("tessedit_char_whitelist", `0123456789xX`) //idCard Must
		defer IDCardTess.Clear()

		mw.SetImageFormat("JPEG")
		blob := mw.GetImageBlob()
		pix, err := leptonica.NewPixReadMem(&blob)

		if err != nil {
			log.Println(err)
			ctx.JSON(200, map[string]interface{}{"error": err})
		}

		IDCardTess.SetImagePix(pix)
		out := IDCardTess.Text()
		log.Println(out)
		//log.Println(IDCardTess.BoxText(0))
		ctx.JSON(200, map[string]interface{}{
			"data": strings.Replace(strings.TrimSpace(out), " ", "", -1),
		})
	}
}
