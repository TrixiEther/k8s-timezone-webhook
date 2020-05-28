NAME = timezone-webhook
OWNER = me
MOD = timezone-webhook

app: deps
	go build -v -o $(NAME) src/webhook.go

deps: mod
	go get -v ./...

mod:
	go mod init github.com/$(OWNER)/$(MOD)