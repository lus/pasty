# Build the application
FROM golang:1.16-alpine AS build
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh build-base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build \
        -o pasty \
        -ldflags "\
            -X github.com/lus/pasty/internal/static.Version=$(git rev-parse --abbrev-ref HEAD)-$(git describe --tags --abbrev=0)-$(git log --pretty=format:'%h' -n 1)" \
        ./cmd/pasty/main.go

# Run the application in an empty alpine environment
FROM alpine:latest
WORKDIR /root
COPY --from=build /app/pasty .
COPY web ./web/
CMD ["./pasty"]