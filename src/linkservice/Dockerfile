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
ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o /cmd/api ./cmd/api/api.go" -command="/cmd/api"

################ Build ################
FROM golang:1.16-alpine as build

COPY --from=modules /go/pkg /go/pkg

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main ./cmd/api/api.go

################ Production ################
FROM gcr.io/distroless/base-debian10 as production
COPY --from=build /app/main /
CMD ["/main"]
