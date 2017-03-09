.PHONY: 

all: tesseract leptonica imagemagick web

web:
	npm install bower
	bower install jquery
	bower install https://git.ishopex.cn/lukaijie/h5_imgcrop.git

tesseract:
	sudo apt-get install libtesseract3 libtesseract-dev
	go get gopkg.in/GeertJohan/go.tesseract.v1

leptonica:
	sudo apt-get install libleptonica-dev
	go get gopkg.in/GeertJohan/go.leptonica.v1

imagemagick:
	sudo apt-get install libmagickwand-dev
	pkg-config --cflags --libs MagickWand
	go get gopkg.in/gographics/imagick.v2/imagick

start:
	@go run main.go server --config .config.yaml
