################ Modules ################
FROM golang:1.16 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

################ Build ################
FROM golang:1.16-buster as build

COPY --from=modules /go/pkg /go/pkg

WORKDIR /app
COPY . .

RUN useradd -u 10001 frontend

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/frontend main.go

################ Production ################
FROM alpine as production

COPY --from=build /etc/passwd /etc/passwd
USER frontend

WORKDIR /frontend
COPY --from=build /go/bin/frontend /frontend/server
COPY ./templates ./templates
COPY ./static ./static

ENTRYPOINT ["/frontend/server"]