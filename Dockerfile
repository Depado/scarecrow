# Build Step
FROM golang:1.16.2-alpine3.13 AS builder

# Dependencies
RUN apk update && apk add --no-cache upx

# Source
WORKDIR $GOPATH/src/github.com/Depado/scarecrow
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
ARG build
ARG version
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=${version} -X main.Build=${build}" -o /tmp/scarecrow
RUN upx /tmp/scarecrow


# Final Step
FROM gcr.io/distroless/static
COPY --from=builder /tmp/scarecrow /go/bin/scarecrow

VOLUME [ "/data" ]
WORKDIR /data
ENTRYPOINT ["/go/bin/scarecrow"]
