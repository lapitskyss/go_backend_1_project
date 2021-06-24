################ Modules ################
FROM golang:1.16 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

################ Develop ################
FROM golang:1.16-buster as develop

COPY --from=modules /go/pkg /go/pkg

ENV CGO_ENABLED=0

WORKDIR /app
COPY . .

RUN go get github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o /cmd/client main.go" -command="/cmd/client"
