set GOOS=linux

set GOARCH=amd64

set CGO_ENABLED=0

go build -ldflags -s -a -installsuffix cgo amqpserver.go