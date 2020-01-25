
AWS_PROFILE:=default
MAKE_ARGS=--profile $(AWS_PROFILE) --region $(AWS_REGION)

SAM_SOURCE_BUCKET:=sam-lambda-artifacts-kp
STACK_NAME:=book-api

build-binary:
	cd src && go mod download && \
	GOOS=linux GOARCH=amd64 go build -o main *.go
package:
	sam package --template-file build/saml.yaml --output-template-file build/packaged.yaml --s3-bucket $(SAM_SOURCE_BUCKET) $(MAKE_ARGS)
deploy: build-binary package
	sam deploy --template-file build/packaged.yaml --stack-name $(STACK_NAME) --capabilities CAPABILITY_IAM $(MAKE_ARGS) 
deploy-only:
	sam deploy --template-file build/packaged.yaml --stack-name $(STACK_NAME) --capabilities CAPABILITY_IAM $(MAKE_ARGS) --no-fail-on-empty-changeset

delete:
	aws cloudformation delete-stack --stack-name $(STACK_NAME) $(MAKE_ARGS)
test:
	echo "Please add tests for this lambda function!"
