build: clean
	GOOS=linux go build main.go
	zip main.zip ./main

start: build
	aws-sam-local validate
	aws-sam-local local start-api 

package: build
	aws cloudformation package --template-file ./sam.yaml --s3-bucket keymesh-tgbot --output-template-file ./packaged.yaml

deploy: package
	aws cloudformation deploy --template-file ./packaged.yaml --stack-name keymesh-tgbot --capabilities CAPABILITY_NAMED_IAM

describe-stack-events:
	aws cloudformation describe-stack-events --stack-name keymesh-tgbot


clean:
	rm -rf main.zip
