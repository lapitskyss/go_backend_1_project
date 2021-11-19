FROM golangci/golangci-lint:v1.41-alpine

WORKDIR /app
COPY . .

RUN golangci-lint run --issues-exit-code=1  ./...
