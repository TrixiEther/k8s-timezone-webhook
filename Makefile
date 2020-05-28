NAME = timezone-webhook
OWNER = me
MOD = timezone-webhook

deploy: docker
	docker run -d -p 8080:8080 $(NAME):latest

docker: app
	docker build -t $(NAME) .

app: deps
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o timezone-webhook ./src/webhook.go

deps: mod
	go get -v ./...

mod: clean
	go mod init github.com/$(OWNER)/$(MOD)

clean:
	rm -f go.mod
	rm -f go.sum
	rm -f $(NAME)
