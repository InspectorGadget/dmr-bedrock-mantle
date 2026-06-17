.PHONY: pull-model run-local run-bedrock build tidy \
        demo-cfn-local demo-cfn-bedrock \
        demo-no-stream-local demo-no-stream-bedrock \
        demo-model help

## pull-model: pull Gemma 3 (one-time, ~2.5 GB)
pull-model:
	docker model pull ai/gemma3

tidy:
	go mod tidy

build:
	go build -o summarizer .

run-local:
	go run . -backend local -file sample-compose.yml

run-bedrock:
	go run . -backend bedrock -file sample-compose.yml

## demo-cfn-local: same thing but with a CloudFormation template
demo-cfn-local:
	go run . -backend local -file sample-cfn.yaml

demo-cfn-bedrock:
	go run . -backend bedrock -file sample-cfn.yaml

## demo-no-stream-local: wait for full response instead of streaming
demo-no-stream-local:
	go run . -backend local -file sample-compose.yml -stream=false

demo-no-stream-bedrock:
	go run . -backend bedrock -file sample-compose.yml -stream=false

## demo-phi4: swap in Phi-4 (pull it first: docker model pull ai/phi4)
demo-phi4:
	go run . -backend local -file sample-compose.yml -model ai/phi4

## demo-llama: swap in Llama 3.2
demo-llama:
	go run . -backend local -file sample-compose.yml -model ai/llama3.2

## run-custom: make run-custom FILE=path/to/file
run-custom:
	go run . -backend local -file $(FILE)

## run-custom-bedrock: make run-custom-bedrock FILE=path/to/file
run-custom-bedrock:
	go run . -backend bedrock -file $(FILE)

## help: list targets
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'
