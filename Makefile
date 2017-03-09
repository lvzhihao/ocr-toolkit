.PHONY: 

web:
	npm install bower
	bower install jquery
	bower install https://git.ishopex.cn/lukaijie/h5_imgcrop.git

start:
	@go run main.go server --config .config.yaml

all: npm
