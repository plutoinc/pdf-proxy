.PHONY: build clean deploy deploy-prod

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/proxy proxy/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	npx sls deploy --verbose --aws-profile pluto-lambda-maker

deploy-prod: clean build
	npx sls deploy -s prod --verbose --aws-profile pluto-lambda-maker
