FROM golang:1.15.2-buster

ENV CGO_ENABLED=0

WORKDIR /app
COPY . .

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o /cmd/client main.go" -command="/cmd/client"
