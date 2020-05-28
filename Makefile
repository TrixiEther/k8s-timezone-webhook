NAME = timezone-webhook

app: deps
	go build -v -o $(NAME) src/webhook.go

deps:
	go get -v ./...