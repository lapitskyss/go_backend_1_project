################ Modules ################
FROM golang:1.17 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

################ Build ################
FROM golang:1.17-buster as build

COPY --from=modules /go/pkg /go/pkg

WORKDIR /app
COPY . .

RUN useradd -u 10001 frontend

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main /app/cmd/main.go

################ Production ################
FROM gcr.io/distroless/base-debian11 as production

COPY --from=build /etc/passwd /etc/passwd
USER frontend

COPY --from=build /app/main /
CMD ["/main"]
