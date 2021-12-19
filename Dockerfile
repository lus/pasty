# Choose the golang image as the build base image
FROM golang:1.16-alpine AS build

# Define the directory we should work in
WORKDIR /app

# Download the necessary go modules
COPY go.mod go.sum ./
RUN go mod download

# Build the application
ARG PASTY_VERSION=unset-debug
COPY . .
RUN go build \
        -o pasty \
        -ldflags "\
            -X github.com/lus/pasty/internal/static.Version=$PASTY_VERSION" \
        ./cmd/pasty/main.go

# Run the application in an empty alpine environment
FROM alpine:latest
WORKDIR /root
COPY --from=build /app/pasty .
COPY web ./web/
EXPOSE 8080
CMD ["./pasty"]