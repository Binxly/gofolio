CSS_INPUT = ./static/assets/css/input.css
CSS_OUTPUT = ./static/assets/css/output.css
LAMBDA_TAGS = lambda,lambda.norpc
FUNCTION_NAME = function

.PHONY: dev watch build clean generate run help css lambda-build deploy quick-deploy

help:
	@echo "Available commands:"
	@echo "  make dev          - Start development server (auto restart)"
	@echo "  make run          - Generate templates and run server once"
	@echo "  make css          - Build CSS once"
	@echo "  make build        - Build production binary"
	@echo "  make lambda-build - Build for AWS Lambda deployment (ARM64)"
	@echo "  make deploy       - Deploy to existing Lambda function"
	@echo "  make quick-deploy - Quick deploy using script"
	@echo "  make clean        - Clean generated files"


generate:
	go tool templ generate

templ-watch:
	go tool templ generate --watch --proxy="http://localhost:3000" --open-browser=true

css:
	tailwindcss -i $(CSS_INPUT) -o $(CSS_OUTPUT)

css-watch:
	tailwindcss -i $(CSS_INPUT) -o $(CSS_OUTPUT) --watch

prepare: generate css

run: prepare
	go run .

build: prepare
	go build -o portfolio .

lambda-build: prepare
	GOOS=linux GOARCH=arm64 go build -tags $(LAMBDA_TAGS) -o bootstrap .
	zip -r function.zip bootstrap static/ assets/ content/

deploy: lambda-build
	aws lambda update-function-code \
		--function-name $(FUNCTION_NAME) \
		--zip-file fileb://function.zip

quick-deploy:
	./scripts/deploy-lambda.sh

server:
	/home/binx/go/bin/air \
	--build.cmd "go build -o tmp/bin/main ." \
	--build.bin "tmp/bin/main" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

dev:
	make -j3 css-watch templ-watch server

clean:
	rm -rf tmp/
	rm -rf bin/
	rm -f *_templ.go
	rm -f bootstrap
	rm -f function.zip