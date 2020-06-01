NAME = timezone-webhook
OWNER = me
MOD = timezone-webhook
VERSION = v1.0

deploy: docker
	docker run -d -p 8080:8080 $(NAME):$(VERSION)

docker: app
	docker build -t $(NAME):$(VERSION) .

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
