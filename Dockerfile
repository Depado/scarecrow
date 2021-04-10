# Build Step
FROM golang:1.16.2-alpine3.13 AS builder

# Dependencies
RUN apk update && apk add --no-cache upx make git

# Source
WORKDIR $GOPATH/src/github.com/Depado/scarecrow
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
RUN make tmp
RUN upx /tmp/scarecrow


# Final Step
FROM gcr.io/distroless/static
COPY --from=builder /tmp/scarecrow /go/bin/scarecrow
ENTRYPOINT ["/go/bin/scarecrow"]
